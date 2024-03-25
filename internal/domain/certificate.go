package domain

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
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
