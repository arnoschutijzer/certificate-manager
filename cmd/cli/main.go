package main

import (
	"time"

	"github.com/algleymi/certificate-manager/internal/exporters"
	"github.com/algleymi/certificate-manager/internal/onepassword"
)

func main() {
	o, err := onepassword.NewOnePassword()

	if err != nil {
		panic(err)
	}

	inAMonth := time.Now().AddDate(0, 1, 0)

	outdatedCertificates, err := o.FindCertificatesOlder(inAMonth)

	if err != nil {
		panic(err)
	}

	exporter := &exporters.ConsoleExporter{}
	exporter.Export(outdatedCertificates)
}
