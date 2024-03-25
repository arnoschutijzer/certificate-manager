package caches

import (
	"time"

	"github.com/algleymi/certificate-manager/internal"
)

type Cache interface {
	SaveVaultItem(vaultItem CachedItem) error
	RetrieveVaultItem(id string) (CachedItem, error)
	UpdateVaultItem(vaultItem CachedItem) error
	Cleanup() error
}

// Decouple vaults from db schema
type CachedItem struct {
	VaultId      string `gorm:"primaryKey"`
	Title        string
	UpdatedAt    time.Time
	Certificates []CachedCertificate `gorm:"foreignKey:VaultItem"`
}

type CachedCertificate struct {
	VaultItem   string
	NotBefore   time.Time
	NotAfter    time.Time
	Subject     string
	CustomName  string
	Fingerprint string `gorm:"primaryKey"`
}

func ToDbCertificate(id string, certificate internal.Certificate) CachedCertificate {
	return CachedCertificate{
		VaultItem:   id,
		Fingerprint: certificate.Fingerprint,
		Subject:     certificate.Subject,
		NotAfter:    certificate.NotAfter,
		NotBefore:   certificate.NotBefore,
		CustomName:  certificate.CustomName,
	}
}

func ToDomainCertificate(certificate CachedCertificate) internal.Certificate {
	return internal.Certificate{
		Fingerprint: certificate.Fingerprint,
		Subject:     certificate.Subject,
		NotAfter:    certificate.NotAfter,
		NotBefore:   certificate.NotBefore,
		CustomName:  certificate.CustomName,
	}
}
