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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPluginhelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pluginhelper Suite")
}

type certs struct {
	serverCertPEM, serverKeyPEM, clientCertPEM []byte
}

func generateCerts(organization []string, serverCommonName string, clientCommonName string) (*certs, error) {
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate server key: %w", err)
	}

	serverCertTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: organization,
			CommonName:   serverCommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	serverCertDER, err := x509.CreateCertificate(
		rand.Reader,
		serverCertTemplate,
		serverCertTemplate,
		&serverKey.PublicKey,
		serverKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create server certificate: %w", err)
	}

	// Generate client key
	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client key: %w", err)
	}

	// Generate client certificate
	clientCertTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: organization,
			CommonName:   clientCommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	clientCertDER, err := x509.CreateCertificate(
		rand.Reader,
		clientCertTemplate,
		serverCertTemplate,
		&clientKey.PublicKey,
		serverKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client certificate: %w", err)
	}

	// Encode server certificate and key to PEM
	serverCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertDER,
	})
	serverPrivateKy, err := x509.MarshalECPrivateKey(serverKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal server private key: %w", err)
	}
	serverKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: serverPrivateKy,
	})

	// Encode client certificate to PEM
	clientCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCertDER,
	})

	return &certs{
		serverCertPEM: serverCertPEM,
		serverKeyPEM:  serverKeyPEM,
		clientCertPEM: clientCertPEM,
	}, nil
}
