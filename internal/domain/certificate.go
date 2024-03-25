package domain

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"strings"
	"time"
)

type Certificate struct {
	NotBefore   time.Time
	NotAfter    time.Time
	Subject     string
	CustomName  string
	Fingerprint string
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

	hash := sha1.Sum(cert.Raw)
	fingerprint := hex.EncodeToString(hash[:]) // 40 characters

	return Certificate{
		NotAfter:    notAfter,
		NotBefore:   notBefore,
		Subject:     subject,
		CustomName:  customName,
		Fingerprint: fingerprint,
	}
}

const (
	START_CERTIFICATE_TEMPLATE = "-----BEGIN CERTIFICATE-----"
	END_CERTIFICATE_TEMPLATE   = "-----END CERTIFICATE-----"
)

func GetCertificatesFromString(secret string, vaultName string) []Certificate {
	// Some secrets are oddly escaped due to them being
	// copy-pasted to Intellij. These are saved with \n directly in the string.
	secret = strings.ReplaceAll(secret, "\\n", "\n")

	certificates := []Certificate{}
	remainingSubstrings := secret

	for true {
		indexOfBegin := strings.Index(remainingSubstrings, START_CERTIFICATE_TEMPLATE)
		indexOfEnd := strings.Index(remainingSubstrings, END_CERTIFICATE_TEMPLATE)

		if indexOfBegin == -1 {
			return certificates
		}

		certificate := remainingSubstrings[indexOfBegin : indexOfEnd+len(END_CERTIFICATE_TEMPLATE)]

		certificates = append(certificates, NewCertificate(certificate, vaultName))
		remainingSubstrings = remainingSubstrings[indexOfEnd+len(END_CERTIFICATE_TEMPLATE):]
	}

	return certificates
}
