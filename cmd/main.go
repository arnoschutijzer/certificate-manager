package main

import (
	"github.com/algleymi/certificate-manager/internal"
)

func main() {
	GetThem(internal.NewOnePasswordStore())
}

func GetThem(store internal.Store) {
	store.FindCertificatesThatAreOutdated()
}
