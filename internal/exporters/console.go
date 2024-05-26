package exporters

import (
	"fmt"

	"github.com/algleymi/certificate-manager/internal"
)

type ConsoleExporter struct{}

var _ Exporter = &ConsoleExporter{}

func (c *ConsoleExporter) Export(certificates []internal.Certificate) error {
	if len(certificates) == 0 {
		fmt.Println("no outdated certificates found")
		return nil
	}

	fmt.Println("outdated certificates found!")
	for _, v := range certificates {
		fmt.Printf("%s (%s) expires at %s\n", v.Subject, v.CustomName, v.NotAfter)
	}
	return nil
}
