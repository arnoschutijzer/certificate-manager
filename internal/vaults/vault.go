package vaults

import (
	"time"

	"github.com/algleymi/certificate-manager/internal"
)

type Vault interface {
	FindCertificatesThatAreOutdated() ([]internal.Certificate, error)
	FindCertificatesOlderThanDate(time.Time) ([]internal.Certificate, error)
}
