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
	"crypto/tls"
	"os"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("BuildTLSConfig", func() {
	var (
		serverCertPath string
		serverKeyPath  string
		clientCertPath string
	)

	writeTempFile := func(data []byte) (string, error) {
		file, err := os.CreateTemp("", "certfile")
		if err != nil {
			return "", err
		}
		defer func() {
			_ = file.Close()
		}()
		_, err = file.Write(data)
		if err != nil {
			return "", err
		}
		return file.Name(), nil
	}

	ginkgo.BeforeEach(func() {
		certs, err := generateCerts([]string{"Test Organization"},
			"localhost",
			"client",
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		serverCertPath, err = writeTempFile(certs.serverCertPEM)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		serverKeyPath, err = writeTempFile(certs.serverKeyPEM)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		clientCertPath, err = writeTempFile(certs.clientCertPEM)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.AfterEach(func() {
		gomega.Expect(os.Remove(serverCertPath)).ToNot(gomega.HaveOccurred())
		gomega.Expect(os.Remove(serverKeyPath)).ToNot(gomega.HaveOccurred())
		gomega.Expect(os.Remove(clientCertPath)).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("should successfully create a TLS config", func(ctx ginkgo.SpecContext) {
		tlsConfig, err := buildTLSConfig(ctx, serverCertPath, serverKeyPath, clientCertPath)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(tlsConfig).ToNot(gomega.BeNil())
		gomega.Expect(tlsConfig.Certificates).To(gomega.HaveLen(1))
		gomega.Expect(tlsConfig.ClientCAs.Subjects()).ToNot(gomega.BeEmpty()) // nolint:staticcheck
		gomega.Expect(tlsConfig.MinVersion).To(gomega.Equal(uint16(tls.VersionTLS13)))
	})
})
