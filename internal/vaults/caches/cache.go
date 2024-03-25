package caches

import (
	"time"

	"github.com/algleymi/certificate-manager/internal"
)

type Cache interface {
	SaveVaultItem(vaultItem Secret) error
	RetrieveVaultItem(id string) (Secret, error)
	UpdateVaultItem(vaultItem Secret) error
	Cleanup() error
}

// Decouple vaults from db schema
type Secret struct {
	Id           string `gorm:"primaryKey"`
	Title        string
	UpdatedAt    time.Time
	Certificates []Certificate `gorm:"foreignKey:VaultItem"`
}

type Certificate struct {
	VaultItem   string
	NotBefore   time.Time
	NotAfter    time.Time
	Subject     string
	CustomName  string
	Fingerprint string `gorm:"primaryKey"`
}

func ToDbCertificate(id string, certificate internal.Certificate) Certificate {
	return Certificate{
		VaultItem:   id,
		Fingerprint: certificate.Fingerprint,
		Subject:     certificate.Subject,
		NotAfter:    certificate.NotAfter,
		NotBefore:   certificate.NotBefore,
		CustomName:  certificate.CustomName,
	}
}

func ToDomainCertificate(certificate Certificate) internal.Certificate {
	return internal.Certificate{
		Fingerprint: certificate.Fingerprint,
		Subject:     certificate.Subject,
		NotAfter:    certificate.NotAfter,
		NotBefore:   certificate.NotBefore,
		CustomName:  certificate.CustomName,
	}
}
