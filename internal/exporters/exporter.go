package exporters

import "github.com/algleymi/certificate-manager/internal"

type Exporter interface {
	Export(certificates []internal.Certificate) error
}

func CreateExporter(format string) Exporter {
	if format == "csv" {
		return &CsvExporter{}
	}

	return &ConsoleExporter{}
}
