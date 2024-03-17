package internal

import (
	"crypto/x509"
	"encoding/pem"
	"strings"
	"time"
)

type Certificate struct {
	NotBefore  time.Time
	NotAfter   time.Time
	Subject    string
	CustomName string
}

func (c *Certificate) IsValid(date time.Time) bool {
	notAfter := c.NotAfter
	notBefore := c.NotBefore

	return notBefore.Before(date) && notAfter.After(date)
}

func NewCertificate(certificate string, customName string) Certificate {
	block, _ := pem.Decode([]byte(certificate))
	cert, _ := x509.ParseCertificate(block.Bytes)

	notAfter := cert.NotAfter
	notBefore := cert.NotBefore
	subject := cert.Subject.CommonName

	return Certificate{
		NotAfter:   notAfter,
		NotBefore:  notBefore,
		Subject:    subject,
		CustomName: customName,
	}
}

func GetCertificatesFromString(secret string, vaultName string) []Certificate {
	// Some secrets are oddly escaped due to them being
	// copy-pasted to Intellij. These are saved with \n directly in the string.
	secret = strings.ReplaceAll(secret, "\\n", "\n")

	certificates := []Certificate{}
	remainingSubstrings := secret

	startCertificateTemplate := "-----BEGIN CERTIFICATE-----"
	endCertificateTemplate := "-----END CERTIFICATE-----"

	for true {
		indexOfBegin := strings.Index(remainingSubstrings, startCertificateTemplate)
		indexOfEnd := strings.Index(remainingSubstrings, endCertificateTemplate)

		if indexOfBegin == -1 {
			return certificates
		}

		certificate := remainingSubstrings[indexOfBegin : indexOfEnd+len(endCertificateTemplate)]

		certificates = append(certificates, NewCertificate(certificate, vaultName))
		remainingSubstrings = remainingSubstrings[indexOfEnd+len(endCertificateTemplate):]
	}

	return certificates
}
