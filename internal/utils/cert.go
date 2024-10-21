package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type Certificate struct {
	Expiry     float64
	CommonName string
}

func DecodeCert(rawCert []byte) (Certificate, error) {
	block, _ := pem.Decode(rawCert)
	if block == nil || block.Type != "CERTIFICATE" {
		return Certificate{}, fmt.Errorf("failed to decode PEM block containing the certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return Certificate{}, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return Certificate{
		Expiry:     float64(cert.NotAfter.Unix()),
		CommonName: cert.Subject.CommonName,
	}, nil
}
