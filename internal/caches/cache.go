package caches

import (
	"time"

	"github.com/algleymi/certificate-manager/internal/domain"
)

type Cache interface {
	SaveSecret(secret Secret) error
	RetrieveSecret(id string) (Secret, error)
	UpdateSecret(secret Secret) error
	Cleanup() error
}

type Secret struct {
	Id           string `gorm:"primaryKey"`
	Title        string
	UpdatedAt    time.Time
	Certificates []Certificate `gorm:"foreignKey:SecretId"`
}

type Certificate struct {
	SecretId    string
	NotBefore   time.Time
	NotAfter    time.Time
	Subject     string
	CustomName  string
	Fingerprint string `gorm:"primaryKey"`
}

func ToDbCertificate(id string, certificate domain.Certificate) Certificate {
	return Certificate{
		SecretId:    id,
		Fingerprint: certificate.Fingerprint,
		Subject:     certificate.Subject,
		NotAfter:    certificate.NotAfter,
		NotBefore:   certificate.NotBefore,
		CustomName:  certificate.CustomName,
	}
}

func ToDomainCertificate(certificate Certificate) domain.Certificate {
	return domain.Certificate{
		Fingerprint: certificate.Fingerprint,
		Subject:     certificate.Subject,
		NotAfter:    certificate.NotAfter,
		NotBefore:   certificate.NotBefore,
		CustomName:  certificate.CustomName,
	}
}
