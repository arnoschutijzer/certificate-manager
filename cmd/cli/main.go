package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/algleymi/certificate-manager/internal/exporters"
	"github.com/algleymi/certificate-manager/internal/onepassword"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s <format>\n", os.Args[0])
		flag.PrintDefaults()
	}
	format := flag.String("format", "text", "switch between \"text\" and \"csv\"")
	flag.Parse()

	o, err := onepassword.NewOnePassword()
	if err != nil {
		panic(err)
	}

	inAMonth := time.Now().AddDate(0, 1, 0)
	outdatedCertificates, err := o.FindCertificatesOlder(inAMonth)
	if err != nil {
		panic(err)
	}

	exporter := exporters.CreateExporter(*format)
	err = exporter.Export(outdatedCertificates)

	if err != nil {
		panic(err)
	}
}
