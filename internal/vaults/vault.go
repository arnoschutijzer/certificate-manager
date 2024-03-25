package vaults

import (
	"time"

	"github.com/algleymi/certificate-manager/internal/domain"
)

type Vault interface {
	FindCertificatesThatAreOutdated() ([]domain.Certificate, error)
	FindCertificatesOlderThanDate(time.Time) ([]domain.Certificate, error)
}
