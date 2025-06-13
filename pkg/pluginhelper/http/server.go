/*
Copyright The CloudNativePG Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"

	"github.com/cloudnative-pg/cnpg-i/pkg/identity"
	"github.com/cloudnative-pg/machinery/pkg/log"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const (
	unixNetwork = "unix"
	tcpNetwork  = "tcp"

	defaultPluginPath = "/plugins"
)

var (
	errNoServerCert = errors.New("TCP server active, but no server-cert value passed")
	errNoServerKey  = errors.New("TCP server active, but no server-key value passed")
	errNoClientCert = errors.New("TCP server active, but no client-cert value passed")
)

// ServerEnricher is the type of functions that can add register
// service implementations in a GRPC server.
type ServerEnricher func(*grpc.Server) error

// CreateMainCmd creates a command to be used as the server side
// for the CNPG-I infrastructure.
func CreateMainCmd(identityImpl identity.IdentityServer, enrichers ...ServerEnricher) *cobra.Command {
	cmd := &cobra.Command{
		Use: "serve",
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			_, err := logr.FromContext(cmd.Context())
			if err != nil {
				// caller did not supply a logger, inject one
				flags := log.NewFlags(zap.Options{Development: viper.GetBool("debug")})
				flags.ConfigureLogging()

				ctx := log.IntoContext(cmd.Context(), log.WithName("cmd_serve"))
				cmd.SetContext(ctx)
			}
		},
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			srv := &Server{
				IdentityImpl:   identityImpl,
				Enrichers:      enrichers,
				ServerCertPath: viper.GetString("server-cert"),
				ServerKeyPath:  viper.GetString("server-key"),
				ClientCertPath: viper.GetString("client-cert"),
				ServerAddress:  viper.GetString("server-address"),
				PluginPath:     viper.GetString("plugin-path"),
			}

			return srv.Start(cmd.Context())
		},
	}

	cmd.PersistentFlags().Bool(
		"debug",
		true,
		"Enable debugging mode, intended to be used only during development",
	)
	_ = viper.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug"))

	cmd.Flags().String(
		"plugin-path",
		"",
		"The plugins socket path",
	)
	_ = viper.BindPFlag("plugin-path", cmd.Flags().Lookup("plugin-path"))

	cmd.Flags().String(
		"server-cert",
		"",
		"The public key to be used for the server process",
	)
	_ = viper.BindPFlag("server-cert", cmd.Flags().Lookup("server-cert"))

	cmd.Flags().String(
		"server-key",
		"",
		"The key to be used for the server process",
	)
	_ = viper.BindPFlag("server-key", cmd.Flags().Lookup("server-key"))

	cmd.Flags().String(
		"client-cert",
		"",
		"The client public key to verify the connection",
	)
	_ = viper.BindPFlag("client-cert", cmd.Flags().Lookup("client-cert"))

	cmd.Flags().String(
		"server-address",
		"",
		"The address where to listen (i.e. 0:9090)",
	)
	_ = viper.BindPFlag("server-address", cmd.Flags().Lookup("server-address"))

	cmd.MarkFlagsRequiredTogether("server-cert", "server-key", "client-cert", "server-address")
	cmd.MarkFlagsMutuallyExclusive("server-cert", "plugin-path")

	return cmd
}

// Server is the main structure to start a GRPC server.
type Server struct {
	IdentityImpl   identity.IdentityServer
	Enrichers      []ServerEnricher
	ServerCertPath string
	ServerKeyPath  string
	ClientCertPath string
	// mutually exclusive with pluginPath
	ServerAddress string
	// mutually exclusive with serverAddress
	PluginPath string
}

// Start starts the server.
func (s *Server) Start(ctx context.Context) error {
	logger := log.FromContext(ctx)

	identityResponse, err := s.IdentityImpl.GetPluginMetadata(
		ctx,
		&identity.GetPluginMetadataRequest{})
	if err != nil {
		logger.Error(err, "Error while querying the identity service")
		return fmt.Errorf("error while querying the identity service: %w", err)
	}

	// Start accepting connections on the socket
	listener, err := s.createListener(ctx, identityResponse)
	if err != nil {
		logger.Error(err, "While starting server")
		return fmt.Errorf("cannot listen on the socket: %w", err)
	}

	// Create GRPC server
	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logFailedRequestsUnaryServerInterceptor(logger),
			loggingUnaryServerInterceptor(logger),
			recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			logFailedRequestsStreamServerInterceptor(logger),
			loggingStreamServerInterceptor(logger),
			recovery.StreamServerInterceptor(),
		),
	}
	if s.isTLSEnabled() {
		certificatesOptions, err := s.setupTLSCerts(ctx)
		if err != nil {
			logger.Error(err, "While setting up TLS authentication")
			return err
		}

		serverOptions = append(serverOptions, *certificatesOptions)
	} else {
		logger.Info("TCP server not active, skipping TLSCerts generation")
	}

	grpcServer := grpc.NewServer(serverOptions...)
	identity.RegisterIdentityServer(
		grpcServer,
		s.IdentityImpl)
	for _, enrich := range s.Enrichers {
		if enrichErr := enrich(grpcServer); enrichErr != nil {
			return enrichErr
		}
	}

	pluginName := identityResponse.GetName()
	pluginDisplayName := identityResponse.GetDisplayName()
	pluginVersion := identityResponse.GetVersion()

	logger.Info(
		"Starting plugin",
		"name", pluginName,
		"displayName", pluginDisplayName,
		"version", pluginVersion,
	)

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()

	if err = grpcServer.Serve(listener); !errors.Is(err, net.ErrClosed) {
		logger.Error(err, "While terminating server")
	}

	return nil
}

func (s *Server) isTLSEnabled() bool {
	return s.ServerCertPath != "" || s.ServerKeyPath != "" || s.ClientCertPath != ""
}

func (s *Server) setupTLSCerts(ctx context.Context) (*grpc.ServerOption, error) {
	logger := log.FromContext(ctx).WithValues(
		"serverCertPath", s.ServerCertPath,
		"serverKeyPath", s.ServerKeyPath,
		"clientCertPath", s.ClientCertPath,
	)

	if s.ServerCertPath == "" {
		return nil, errNoServerCert
	}

	if s.ServerKeyPath == "" {
		return nil, errNoServerKey
	}

	if s.ClientCertPath == "" {
		return nil, errNoClientCert
	}

	tlsConfig, err := s.buildTLSConfig(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info("Set up TLS authentication")
	result := grpc.Creds(credentials.NewTLS(tlsConfig))

	return &result, nil
}

func (s *Server) buildTLSConfig(ctx context.Context) (*tls.Config, error) {
	logger := log.FromContext(ctx).WithValues(
		"serverCertPath", s.ServerCertPath,
		"serverKeyPath", s.ServerKeyPath,
		"clientCertPath", s.ClientCertPath,
	)

	return &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		GetCertificate: func(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := tls.LoadX509KeyPair(s.ServerCertPath, s.ServerKeyPath)
			if err != nil {
				logger.Error(err, "failed to load server key pair")
				return nil, fmt.Errorf("failed to load server key pair: %w", err)
			}
			return &cert, nil
		},
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			caCertPool := x509.NewCertPool()
			caBytes, err := os.ReadFile(filepath.Clean(s.ClientCertPath))
			if err != nil {
				logger.Error(err, "failed to read client public key")
				return fmt.Errorf("failed to read client public key: %w", err)
			}
			if ok := caCertPool.AppendCertsFromPEM(caBytes); !ok {
				logger.Error(err, "failed to parse client public key")
				return fmt.Errorf("failed to parse client public key: %w", err)
			}
			// Parse the server certificate
			certs := make([]*x509.Certificate, len(rawCerts))
			for i, asn1Data := range rawCerts {
				cert, err := x509.ParseCertificate(asn1Data)
				if err != nil {
					return fmt.Errorf("failed to parse server certificate: %w", err)
				}
				certs[i] = cert
			}
			opts := x509.VerifyOptions{
				Roots: caCertPool,
			}
			_, err = certs[0].Verify(opts)
			return err
		},
		InsecureSkipVerify: true, // Required to use VerifyPeerCertificate
		MinVersion:         tls.VersionTLS13,
	}, nil
}

func (s *Server) createListener(
	ctx context.Context,
	metadata *identity.GetPluginMetadataResponse,
) (net.Listener, error) {
	if len(s.ServerAddress) != 0 {
		return s.createTCPListener(ctx)
	}

	return s.createUnixDomainSocketListener(ctx, metadata)
}

func (s *Server) createTCPListener(ctx context.Context) (net.Listener, error) {
	logger := log.FromContext(ctx)

	logger.Info(
		"Starting plugin listener",
		"protocol", tcpNetwork,
		"serverAddress", s.ServerAddress,
	)

	// Start accepting connections on the server address
	listener, err := net.Listen(
		tcpNetwork,
		s.ServerAddress,
	)
	if err != nil {
		logger.Error(err, "While starting server")
		return nil, fmt.Errorf("cannot listen on `%s`: %w", s.ServerAddress, err)
	}

	return listener, nil
}

func (s *Server) createUnixDomainSocketListener(
	ctx context.Context,
	metadata *identity.GetPluginMetadataResponse,
) (net.Listener, error) {
	logger := log.FromContext(ctx)

	pluginPath := s.PluginPath
	if len(s.PluginPath) == 0 {
		pluginPath = defaultPluginPath
	}
	socketName := path.Join(pluginPath, metadata.GetName())

	// Remove stale unix socket it still existent
	if err := s.removeStaleSocket(ctx, socketName); err != nil {
		logger.Error(err, "While removing old unix socket")
		return nil, err
	}

	logger.Info(
		"Starting plugin listener",
		"protocol", unixNetwork,
		"socketName", socketName,
	)

	// Start accepting connections on the socket
	listener, err := net.Listen(
		unixNetwork,
		socketName,
	)
	if err != nil {
		logger.Error(err, "While starting server")
		return nil, fmt.Errorf("cannot listen on `%s`: %w", socketName, err)
	}

	return listener, nil
}

// removeStaleSocket removes a stale unix domain socket.
func (s *Server) removeStaleSocket(ctx context.Context, pluginPath string) error {
	logger := log.FromContext(ctx)
	_, err := os.Stat(pluginPath)

	switch {
	case err == nil:
		logger.Info("Removing stale socket", "pluginPath", pluginPath)
		err := os.Remove(pluginPath)
		if err != nil {
			return fmt.Errorf("error while removing stale socket: %w", err)
		}

		return nil

	case errors.Is(err, os.ErrNotExist):
		return nil

	default:
		return fmt.Errorf("error while checking for stale socket: %w", err)
	}
}
