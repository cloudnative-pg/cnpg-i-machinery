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
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cloudnative-pg/machinery/pkg/log"
)

// TLSConfigManager manages dynamic TLS configuration
type TLSConfigManager struct {
	mu            sync.RWMutex
	serverCert    string
	serverKey     string
	clientCAFile  string
	currentConfig *tls.Config
}

// newTLSConfigManager creates a new TLS config manager
func newTLSConfigManager(certFile, keyFile, caFile string) (*TLSConfigManager, error) {
	manager := &TLSConfigManager{
		serverCert:   certFile,
		serverKey:    keyFile,
		clientCAFile: caFile,
	}

	// Load initial configuration
	if err := manager.reload(); err != nil {
		return nil, err
	}

	return manager, nil
}

// reload loads the TLS configuration from files
func (m *TLSConfigManager) reload() error {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair(m.serverCert, m.serverKey)
	if err != nil {
		return fmt.Errorf("failed to load server cert: %v", err)
	}

	// Load client CA
	clientCACert, err := os.ReadFile(filepath.Clean(m.clientCAFile))
	if err != nil {
		return fmt.Errorf("failed to read client CA: %v", err)
	}

	clientCAs := x509.NewCertPool()
	if !clientCAs.AppendCertsFromPEM(clientCACert) {
		return fmt.Errorf("failed to parse client CA")
	}

	// Create new TLS config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
	}

	// Update the configuration atomically
	m.mu.Lock()
	m.currentConfig = tlsConfig
	m.mu.Unlock()

	log.Info("TLS configuration reloaded successfully")
	return nil
}

// GetConfigForConnection returns the current TLS configuration
// This is called for each new connection
func (m *TLSConfigManager) GetConfigForConnection(_ *tls.ClientHelloInfo) (*tls.Config, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentConfig.Clone(), nil
}

// Watch monitors for configuration changes and reloads
func (m *TLSConfigManager) Watch(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.reload(); err != nil {
				log.Error(err, "Failed to reload TLS config")
			}
		}
	}
}
