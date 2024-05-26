package exporters

import "github.com/algleymi/certificate-manager/internal"

type Exporter interface {
	Export(certificates []internal.Certificate) error
}
