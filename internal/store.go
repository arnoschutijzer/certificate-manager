package internal

import "time"

type Store interface {
	FindCertificatesThatAreOutdated() ([]Item, error)
	FindCertificatesOlderThanDate(time.Time) ([]Item, error)
}
