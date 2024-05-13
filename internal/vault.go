package internal

import "time"

type Vault interface {
	FindCertificatesOlder(time.Time) ([]Certificate, error)
}
