package caches

import (
	"time"

	"github.com/algleymi/certificate-manager/internal"
	"gorm.io/gorm"
)

type Cache interface {
	SaveVaultItem(vaultItem VaultItem) error
	RetrieveVaultItem(id string)
	Cleanup() error
}

// Decouple vaults from db schema
type VaultItem struct {
	gorm.Model
	VaultId      string
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
	Fingerprint string
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
