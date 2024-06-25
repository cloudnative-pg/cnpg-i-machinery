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

package logging

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func newLogger(debug bool) logr.Logger {
	var zapLog *zap.Logger
	var err error

	if debug {
		zapLog, err = zap.NewDevelopment()
	} else {
		zapLog, err = zap.NewProduction()
	}
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}

	result := zapr.NewLogger(zapLog)
	return result
}

// IntoContext injects the logger into the passed context, returning
// a context having the logger embedded. The logger can be recovered
// with FromContext
func IntoContext(ctx context.Context, logger logr.Logger) context.Context {
	return logr.NewContext(ctx, logger)
}

// NewIntoContext injects a new logger into the passed context, returning
// a context having the logger embedded. The logger can be recovered
// with FromContext
func NewIntoContext(ctx context.Context, debug bool) context.Context {
	logger := newLogger(debug)
	return IntoContext(ctx, logger)
}

// FromContext get the logger from the context, generating a new generic
// logger if one is not found.
//
// This should probably have a means of panicking if a logger is not found
// during development.
//
func FromContext(ctx context.Context) logr.Logger {
	logger, err := logr.FromContext(ctx)
	if err != nil {
		return newLogger(false)
	}
	return logger
}
