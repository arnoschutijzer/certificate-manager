package internal

import (
	"time"

	"github.com/algleymi/certificate-manager/internal"
)

type Vault interface {
	FindCertificatesThatAreOutdated() ([]internal.Certificate, error)
	FindCertificatesOlderThanDate(time.Time) ([]internal.Certificate, error)
	FindCertificatesOlderThanDate2(time.Time) ([]internal.Certificate, error)
}
