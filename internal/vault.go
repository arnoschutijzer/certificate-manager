package internal

import "time"

type Vault interface {
	FindCertificatesThatAreOutdated() ([]Certificate, error)
	FindCertificatesOlderThanDate(time.Time) ([]Certificate, error)
}
