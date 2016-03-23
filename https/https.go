package https

import (
	"os"

	"crypto/tls"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func ListenAndServeTLS(bindAddress string, tlsCert string, tlsKey string, handler http.Handler) error {
	// Checking for TLS default keys
	var certBytes, keyBytes []byte

	if _, err := os.Stat(tlsCert); err != nil {
		if _, err := os.Stat(tlsKey); err != nil {
			certBytes, keyBytes = GenerateDefaultTLS(tlsCert, tlsKey)
		}

	} else {
		log.Info("Loading TLS Certificate: ", tlsCert)
		log.Info("Loading TLS Private key: ", tlsKey)

		certBytes, err = ioutil.ReadFile(tlsCert)
		keyBytes, err = ioutil.ReadFile(tlsKey)
	}

	certifs, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		log.Panic("Cannot parse certificates")
	}

	s := &http.Server{
		Addr:    bindAddress,
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{certifs},
		},
	}

	log.Info("Listening (https) on ", bindAddress)

	return s.ListenAndServeTLS("", "")
}
