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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"sync"

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
	if err := manager.load(); err != nil {
		return nil, err
	}

	return manager, nil
}

// load loads the TLS configuration from files
func (m *TLSConfigManager) load() error {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair(m.serverCert, m.serverKey)
	if err != nil {
		return fmt.Errorf("failed to load server cert from %s and %s: %w", m.serverCert, m.serverKey, err)
	}

	// Load client CA
	clientCACert, err := os.ReadFile(filepath.Clean(m.clientCAFile))
	if err != nil {
		return fmt.Errorf("failed to read client CA from %s: %w", m.clientCAFile, err)
	}

	clientCAs := x509.NewCertPool()
	if !clientCAs.AppendCertsFromPEM(clientCACert) {
		return fmt.Errorf("failed to parse client CA from %s", m.clientCAFile)
	}

	// Create new TLS config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    clientCAs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}

	// Update the configuration atomically
	m.mu.Lock()
	m.currentConfig = tlsConfig
	m.mu.Unlock()

	log.Info("TLS configuration loaded successfully",
		"serverCert", m.serverCert,
		"serverKey", m.serverKey,
		"clientCA", m.clientCAFile)
	return nil
}

// Reload reloads the TLS configuration from files (public method)
func (m *TLSConfigManager) Reload() error {
	return m.load()
}

// GetConfigForConnection returns the current TLS configuration
// This is called for each new connection
func (m *TLSConfigManager) GetConfigForConnection(_ *tls.ClientHelloInfo) (*tls.Config, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentConfig.Clone(), nil
}
