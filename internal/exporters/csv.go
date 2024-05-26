package exporters

import (
	"fmt"
	"os"

	"github.com/algleymi/certificate-manager/internal"
)

type CsvExporter struct{}

var _ Exporter = &CsvExporter{}

func (c *CsvExporter) Export(certificates []internal.Certificate) error {
	destination, err := os.Create("outdated-certificates.csv")

	if err != nil {
		return err
	}

	fmt.Fprintf(destination, "%s;%s;%s;%s\n", "CustomName", "Subject", "Fingerprint", "NotAfter")

	for _, v := range certificates {
		fmt.Fprintf(destination, "%s;%s;%s;%s\n", v.CustomName, v.Subject, v.Fingerprint, v.NotAfter)
	}
	return nil
}
