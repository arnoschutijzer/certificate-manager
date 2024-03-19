package caches

import (
	"time"

	i "github.com/algleymi/certificate-manager/internal"
)

type Cache interface {
	SaveCertificate(certificate i.Certificate)
	RetrieveCertificate(fingerprint string)
}

// Decouple from db schema
type DatabaseCertificate struct {
	// gorm.Model
	Fingerprint string `gorm:"primaryKey"`
	CustomName  string
	Subject     string
	NotBefore   time.Time
	NotAfter    time.Time
}

func ToDatabaseCertificate(certificate i.Certificate) DatabaseCertificate {
	return DatabaseCertificate{
		Fingerprint: certificate.Fingerprint,
		CustomName:  certificate.CustomName,
		Subject:     certificate.Subject,
		NotBefore:   certificate.NotBefore,
		NotAfter:    certificate.NotAfter,
	}
}
