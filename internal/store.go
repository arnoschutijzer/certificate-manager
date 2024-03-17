package internal

import "time"

type Store interface {
	FindCertificatesThatAreOutdated() ([]Certificate, error)
	FindCertificatesOlderThanDate(time.Time) ([]Certificate, error)
}
