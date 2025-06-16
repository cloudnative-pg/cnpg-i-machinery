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
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BuildTLSConfig", func() {
	var server Server

	writeTempFile := func(data []byte) (string, error) {
		file, err := os.CreateTemp("", "certfile")
		if err != nil {
			return "", fmt.Errorf("failed to create temp file: %w", err)
		}
		defer func() {
			_ = file.Close()
		}()
		_, err = file.Write(data)
		if err != nil {
			return "", fmt.Errorf("failed to write to temp file: %w", err)
		}

		return file.Name(), nil
	}

	BeforeEach(func() {
		certs, err := generateCerts([]string{"Test Organization"},
			"localhost",
			"client",
		)
		Expect(err).ToNot(HaveOccurred())

		serverCertPath, err := writeTempFile(certs.serverCertPEM)
		Expect(err).ToNot(HaveOccurred())

		serverKeyPath, err := writeTempFile(certs.serverKeyPEM)
		Expect(err).ToNot(HaveOccurred())

		clientCertPath, err := writeTempFile(certs.clientCertPEM)
		Expect(err).ToNot(HaveOccurred())

		server = Server{
			IdentityImpl:   nil,
			Enrichers:      nil,
			ServerCertPath: serverCertPath,
			ServerKeyPath:  serverKeyPath,
			ClientCertPath: clientCertPath,
			ServerAddress:  "",
			PluginPath:     "",
		}
	})

	AfterEach(func() {
		Expect(os.Remove(server.ServerCertPath)).ToNot(HaveOccurred())
		Expect(os.Remove(server.ServerKeyPath)).ToNot(HaveOccurred())
		Expect(os.Remove(server.ClientCertPath)).ToNot(HaveOccurred())
	})

	It("should successfully create a TLS config", func(ctx SpecContext) {
		tlsConfigManager, err := newTLSConfigManager(server.ServerCertPath, server.ServerKeyPath, server.ClientCertPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(tlsConfigManager).ToNot(BeNil())
		tlsConfig := tlsConfigManager.currentConfig
		Expect(tlsConfig).ToNot(BeNil())
		Expect(tlsConfig.GetCertificate).ToNot(BeNil())
		Expect(tlsConfig.ClientCAs.Subjects()).ToNot(BeEmpty()) //nolint: staticcheck
		Expect(tlsConfig.MinVersion).To(Equal(uint16(tls.VersionTLS13)))

		cert, err := tlsConfig.GetCertificate(nil)
		Expect(err).Error().NotTo(HaveOccurred())
		Expect(cert).NotTo(BeNil())
	})
})
