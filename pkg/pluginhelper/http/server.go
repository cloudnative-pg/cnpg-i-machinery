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
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

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
			if err == nil {
				// caller did not supply a logger, inject one
				flags := log.NewFlags(zap.Options{Development: viper.GetBool("debug")})
				flags.ConfigureLogging()

				ctx := log.IntoContext(cmd.Context(), log.WithName("cmd_serve"))
				cmd.SetContext(ctx)
			}
		},
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return run(cmd.Context(), identityImpl, enrichers...)
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

// run starts listening for GRPC requests.
func run(ctx context.Context, identityImpl identity.IdentityServer, enrichers ...ServerEnricher) error {
	logger := log.FromContext(ctx)

	identityResponse, err := identityImpl.GetPluginMetadata(
		ctx,
		&identity.GetPluginMetadataRequest{})
	if err != nil {
		logger.Error(err, "Error while querying the identity service")
		return fmt.Errorf("error while querying the identity service: %w", err)
	}

	// Start accepting connections on the socket
	listener, err := createListener(ctx, identityResponse)
	if err != nil {
		logger.Error(err, "While starting server")
		return fmt.Errorf("cannot listen on the socket: %w", err)
	}

	// Handle quit-like signal
	handleSignals(ctx, listener)

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
	if isTLSEnabled() {
		certificatesOptions, err := setupTLSCerts(ctx)
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
		identityImpl)
	for _, enrich := range enrichers {
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
	if err = grpcServer.Serve(listener); !errors.Is(err, net.ErrClosed) {
		logger.Error(err, "While terminating server")
	}

	return nil
}

func isTLSEnabled() bool {
	serverCertPath := viper.GetString("server-cert")
	serverKeyPath := viper.GetString("server-key")
	clientCertPath := viper.GetString("client-cert")

	return serverCertPath != "" || serverKeyPath != "" || clientCertPath != ""
}

func setupTLSCerts(ctx context.Context) (*grpc.ServerOption, error) {
	serverCertPath := viper.GetString("server-cert")
	serverKeyPath := viper.GetString("server-key")
	clientCertPath := viper.GetString("client-cert")

	logger := log.FromContext(ctx).WithValues(
		"serverCertPath", serverCertPath,
		"serverKeyPath", serverKeyPath,
		"clientCertPath", clientCertPath,
	)

	if serverCertPath == "" {
		return nil, errNoServerCert
	}

	if serverKeyPath == "" {
		return nil, errNoServerKey
	}

	if clientCertPath == "" {
		return nil, errNoClientCert
	}

	tlsConfig, err := buildTLSConfig(ctx, serverCertPath, serverKeyPath, clientCertPath)
	if err != nil {
		return nil, err
	}

	logger.Info("Set up TLS authentication")
	result := grpc.Creds(credentials.NewTLS(tlsConfig))

	return &result, nil
}

func buildTLSConfig(
	ctx context.Context,
	serverCertPath string,
	serverKeyPath string,
	clientCertPath string,
) (*tls.Config, error) {
	logger := log.FromContext(ctx).WithValues(
		"serverCertPath", serverCertPath,
		"serverKeyPath", serverKeyPath,
		"clientCertPath", clientCertPath,
	)

	cert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		logger.Error(err, "failed to load server key pair")
		return nil, fmt.Errorf("failed to load server key pair: %w", err)
	}

	caCertPool := x509.NewCertPool()
	caBytes, err := os.ReadFile(filepath.Clean(clientCertPath))
	if err != nil {
		logger.Error(err, "failed to read client public key")
		return nil, fmt.Errorf("failed to read client public key: %w", err)
	}
	if ok := caCertPool.AppendCertsFromPEM(caBytes); !ok {
		logger.Error(err, "failed to parse client public key")
		return nil, fmt.Errorf("failed to parse client public key: %w", err)
	}

	return &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

func createListener(ctx context.Context, metadata *identity.GetPluginMetadataResponse) (net.Listener, error) {
	serverAddress := viper.GetString("server-address")
	if len(serverAddress) != 0 {
		return createTCPListener(ctx)
	}

	return createUnixDomainSocketListener(ctx, metadata)
}

func createTCPListener(ctx context.Context) (net.Listener, error) {
	logger := log.FromContext(ctx)

	serverAddress := viper.GetString("server-address")

	logger.Info(
		"Starting plugin listener",
		"protocol", tcpNetwork,
		"serverAddress", serverAddress,
	)

	// Start accepting connections on the server address
	listener, err := net.Listen(
		tcpNetwork,
		serverAddress,
	)
	if err != nil {
		logger.Error(err, "While starting server")
		return nil, fmt.Errorf("cannot listen on `%s`: %w", serverAddress, err)
	}

	return listener, nil
}

func createUnixDomainSocketListener(
	ctx context.Context,
	metadata *identity.GetPluginMetadataResponse,
) (net.Listener, error) {
	logger := log.FromContext(ctx)

	pluginPath := viper.GetString("plugin-path")
	if len(pluginPath) == 0 {
		pluginPath = defaultPluginPath
	}
	socketName := path.Join(pluginPath, metadata.GetName())

	// Remove stale unix socket it still existent
	if err := removeStaleSocket(ctx, socketName); err != nil {
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
func removeStaleSocket(ctx context.Context, pluginPath string) error {
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

// handleSignals makes sure that we close the listening socket
// when we receive a quit-like signal.
func handleSignals(ctx context.Context, listener net.Listener) {
	logger := log.FromContext(ctx)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
	go func(c chan os.Signal) {
		sig := <-c
		logger.Info(
			"Caught signal, shutting down.",
			"signal", sig.String())

		if err := listener.Close(); err != nil {
			logger.Error(err, "While stopping server")
		}

		os.Exit(1)
	}(sigc)
}
