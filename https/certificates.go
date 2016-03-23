package https

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"

	"bytes"
	"math/big"
	"time"

	log "github.com/Sirupsen/logrus"
)

// GenerateDefaultTLS will generate a certificate and a private key, on-the-fly
// This certificate will not be persistant (until you save it yourself)
// This is made by design to generate a new certificate at each server start, if you
// don't provide one yourself
func GenerateDefaultTLS(certPath string, keyPath string) (cert []byte, key []byte) {
	log.Debug("Generating TLS certificate and key")

	rootName := "itsyou.online"
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)

	if err != nil {
		log.Panic("Failed to generate private key: ", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(4 * 365 * 24 * time.Hour) // 4 years

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Panic("Failed to generate serial number: ", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:      []string{"Belgium"},
			Locality:     []string{"Lochristi"},
			Organization: []string{"It's You Online"},
			CommonName:   rootName,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,

		DNSNames: []string{rootName},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Panic("Failed to create certificate: ", err)
	}

	// Certificates output (in memory)

	var certOut, keyOut bytes.Buffer

	pem.Encode(&certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		log.Panic("Unable to marshal ECDSA private key: ", err)
	}

	pemUnmarsh := pem.Block{Type: "EC PRIVATE KEY", Bytes: b}

	pem.Encode(&keyOut, &pemUnmarsh)

	return certOut.Bytes(), keyOut.Bytes()
}
