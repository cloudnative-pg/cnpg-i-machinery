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
	"errors"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/cloudnative-pg/cnpg-i/pkg/identity"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/cloudnative-pg/cnpg-i-machinery/pkg/logging"
)

const unixNetwork = "unix"

// ServerEnricher is the type of functions that can add register
// service implementations in a GRPC server
type ServerEnricher func(*grpc.Server) error

// CreateMainCmd creates a command to be used as the server side
// for the CNPG-I infrastructure
func CreateMainCmd(identityImpl identity.IdentityServer, enrichers ...ServerEnricher) *cobra.Command {
	cmd := &cobra.Command{
		Use: "pvc-backup",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := logging.NewIntoContext(
				cmd.Context(),
				viper.GetBool("debug"))
			cmd.SetContext(ctx)
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
		"/plugins",
		"The plugins socket path",
	)
	_ = viper.BindPFlag("plugin-path", cmd.Flags().Lookup("plugin-path"))

	return cmd
}

// run starts listining for GRPC requests
func run(ctx context.Context, identityImpl identity.IdentityServer, enrichers ...ServerEnricher) error {
	logger := logging.FromContext(ctx)

	identityResponse, err := identityImpl.GetPluginMetadata(
		ctx,
		&identity.GetPluginMetadataRequest{})
	if err != nil {
		logger.Error(err, "Error while querying the identity service")
		return err
	}

	pluginPath := viper.GetString("plugin-path")
	pluginName := identityResponse.Name
	pluginDisplayName := identityResponse.DisplayName
	pluginVersion := identityResponse.Version
	socketName := path.Join(pluginPath, identityResponse.Name)

	// Remove stale unix socket it still existent
	if err := removeStaleSocket(ctx, socketName); err != nil {
		logger.Error(err, "While removing old unix socket")
		return err
	}

	// Start accepting connections on the socket
	listener, err := net.Listen(
		unixNetwork,
		socketName,
	)
	if err != nil {
		logger.Error(err, "While starting server")
		return err
	}

	// Handle quit-like signal
	handleSignals(ctx, listener)

	// Create GRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(panicRecoveryHandler(listener))),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(panicRecoveryHandler(listener))),
		),
	)
	identity.RegisterIdentityServer(
		grpcServer,
		identityImpl)
	for _, enrich := range enrichers {
		err := enrich(grpcServer)
		if err != nil {
			return err
		}
	}

	logger.Info(
		"Starting plugin",
		"path", pluginPath,
		"name", pluginName,
		"displayName", pluginDisplayName,
		"version", pluginVersion,
		"socketName", socketName,
	)

	if err = grpcServer.Serve(listener); !errors.Is(err, net.ErrClosed) {
		logger.Error(err, "While terminating server")
	}

	return nil
}

// removeStaleSocket removes a stale unix domain socket
func removeStaleSocket(ctx context.Context, pluginPath string) error {
	logger := logging.FromContext(ctx)
	_, err := os.Stat(pluginPath)

	switch {
	case err == nil:
		logger.Info("Removing stale socket", "pluginPath", pluginPath)
		return os.Remove(pluginPath)

	case errors.Is(err, os.ErrNotExist):
		return nil

	default:
		return err
	}
}

// handleSignals makes sure that we close the listening socket
// when we receive a quit-like signal
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
