package internal

import (
	"crypto/x509"
	"encoding/pem"
	"strings"
	"time"
)

type Certificate struct {
	NotBefore time.Time
	NotAfter  time.Time
	Subject   string
	VaultName string
}

func (c *Certificate) IsValid(date time.Time) bool {
	notAfter := c.NotAfter
	notBefore := c.NotBefore

	return notBefore.Before(date) && notAfter.After(date)
}

func NewCertificate(certificate string, vaultName string) Certificate {
	block, _ := pem.Decode([]byte(certificate))
	cert, _ := x509.ParseCertificate(block.Bytes)

	notAfter := cert.NotAfter
	notBefore := cert.NotBefore
	subject := cert.Subject.CommonName

	return Certificate{
		NotAfter:  notAfter,
		NotBefore: notBefore,
		Subject:   subject,
		VaultName: vaultName,
	}
}

func GetCertificatesFromString(secret string) []string {
	certificates := []string{}
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

		certificates = append(certificates, certificate)
		remainingSubstrings = remainingSubstrings[indexOfEnd+len(endCertificateTemplate):]
	}

	return certificates
}
