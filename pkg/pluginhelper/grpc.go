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

	"github.com/go-logr/logr"
	"google.golang.org/grpc"

	"github.com/cloudnative-pg/cnpg-i-machinery/pkg/logging"
)

// Inject the passed logger into the gRPC call context for all inbound unary calls.
//
// Works around go-grpc's lack of a WithContext option to set a root context.
func loggingUnaryServerInterceptor(logger logr.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		newCtx := logging.IntoContext(ctx, logger)
		return handler(newCtx, req)
	}
}

// Logs failed requests
func logFailedRequestsUnaryServerInterceptor(logger logr.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		result, err := handler(ctx, req)
		if err != nil {
			logger.Error(
				err,
				"Error while handling GRPC request",
				"info", info,
			)
		}
		return result, err
	}
}

// logInjectStream wraps a grpc.ServerStream and injects a logger into the context.
type logInjectStream struct {
	grpc.ServerStream
	logger logr.Logger
}

// Inject the passed logger into the gRPC call context for all inbound streaming calls.
func (s *logInjectStream) Context() context.Context {
	return logging.IntoContext(s.ServerStream.Context(), s.logger)
}

// Inject the passed logger into the gRPC call context for all inbound streaming calls
// by wrapping the ServerStream and overriding the Context() method.
//
// Works around go-grpc's lack of a WithContext option to set a root context.
func loggingStreamServerInterceptor(logger logr.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &logInjectStream{ss, logger})
	}
}

// Logs failed requests
func logFailedRequestsStreamServerInterceptor(logger logr.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err != nil {
			logger.Error(
				err,
				"Error while handling GRPC request",
				"info", info,
			)
		}
		return err
	}
}
