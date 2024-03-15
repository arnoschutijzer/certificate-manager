package internal

type Store interface {
	FindCertificatesThatAreOutdated() []string
}
