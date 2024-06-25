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

package pluginhelper

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
	"syscall"

	"github.com/cloudnative-pg/cnpg-i/pkg/identity"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/cloudnative-pg/cnpg-i-machinery/pkg/logging"
)

const (
	unixNetwork = "unix"
	tcpNetwork  = "tcp"

	defaultPluginPath = "/plugins"
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
				ctx := logging.NewIntoContext(
					cmd.Context(),
					viper.GetBool("debug"))
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
		"Enable debugging mode",
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
		"The key key to be used for the server process",
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
	logger := logging.FromContext(ctx)

	identityResponse, err := identityImpl.GetPluginMetadata(
		ctx,
		&identity.GetPluginMetadataRequest{})
	if err != nil {
		logger.Error(err, "Error while querying the identity service")
		return fmt.Errorf("error while querying the identity service: %w", err)
	}

	pluginName := identityResponse.GetName()
	pluginDisplayName := identityResponse.GetDisplayName()
	pluginVersion := identityResponse.GetVersion()

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
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(panicRecoveryHandler(listener))),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(panicRecoveryHandler(listener))),
		),
	}
	if certificatesOptions, err := setupTLSCerts(ctx); err != nil {
		logger.Error(err, "While setting up TLS authentication")
		return err
	} else if certificatesOptions != nil {
		serverOptions = append(serverOptions, *certificatesOptions)
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

func setupTLSCerts(ctx context.Context) (*grpc.ServerOption, error) {
	serverCertPath := viper.GetString("server-cert")
	serverKeyPath := viper.GetString("server-key")
	clientCertPath := viper.GetString("client-cert")

	// There's no need to load the TLS stuff
	// if the TCP server is not active
	if serverCertPath == "" {
		return nil, nil
	}

	logger := logging.FromContext(ctx).WithValues(
		"serverCertPath", serverCertPath,
		"serverKeyPath", serverKeyPath,
		"clientCertPath", clientCertPath,
	)

	cert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		logger.Error(err, "failed to load server key pair")
		return nil, err
	}

	ca := x509.NewCertPool()
	caBytes, err := os.ReadFile(clientCertPath) // nolint: gosec
	if err != nil {
		logger.Error(err, "failed to read client public key")
		return nil, err
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		logger.Error(err, "failed to parse client public key")
		return nil, err
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
		MinVersion:   tls.VersionTLS13,
	}

	logger.Info("Set up TLS authentication")
	result := grpc.Creds(credentials.NewTLS(tlsConfig))
	return &result, nil
}

func createListener(ctx context.Context, metadata *identity.GetPluginMetadataResponse) (net.Listener, error) {
	serverAddress := viper.GetString("server-address")
	if len(serverAddress) != 0 {
		return createTCPListener(ctx)
	}

	return createUnixDomainSocketListener(ctx, metadata)
}

func createTCPListener(ctx context.Context) (net.Listener, error) {
	logger := logging.FromContext(ctx)

	serverAddress := viper.GetString("server-address")

	logger.Info(
		"Starting plugin listener",
		"protocol", tcpNetwork,
		"serverAddress", serverAddress,
	)

	// Start accepting connections on the socket
	listener, err := net.Listen(
		tcpNetwork,
		serverAddress,
	)
	return listener, err
}

func createUnixDomainSocketListener(
	ctx context.Context,
	metadata *identity.GetPluginMetadataResponse,
) (net.Listener, error) {
	logger := logging.FromContext(ctx)

	pluginPath := viper.GetString("plugin-path")
	if len(pluginPath) == 0 {
		pluginPath = defaultPluginPath
	}
	socketName := path.Join(pluginPath, metadata.Name)

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
	return listener, err
}

// removeStaleSocket removes a stale unix domain socket.
func removeStaleSocket(ctx context.Context, pluginPath string) error {
	logger := logging.FromContext(ctx)
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
	logger := logging.FromContext(ctx)

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

func panicRecoveryHandler(listener net.Listener) recovery.RecoveryHandlerFuncContext {
	return func(ctx context.Context, err any) error {
		logger := logging.FromContext(ctx)
		logger.Info("Panic occurred", "error", err)

		if closeError := listener.Close(); closeError != nil {
			logger.Error(closeError, "While stopping server")
		}

		os.Exit(1)

		return nil
	}
}
