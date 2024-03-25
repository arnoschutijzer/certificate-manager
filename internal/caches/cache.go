package caches

import (
	"time"

	"github.com/algleymi/certificate-manager/internal"
)

type Cache interface {
	SaveVaultItem(vaultItem VaultItem) error
	RetrieveVaultItem(id string) (VaultItem, error)
	UpdateVaultItem(vaultItem VaultItem) error
	Cleanup() error
}

// Decouple vaults from db schema
type VaultItem struct {
	VaultId      string `gorm:"primaryKey"`
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
