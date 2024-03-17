package internal

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParsesASingleCertificate(t *testing.T) {
	certificateAsString, err := os.ReadFile("./test_fixtures/RootCA.pem")
	if err != nil {
		t.Fail()
	}

	firstOfJune2024UTCAt1AM := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)

	certificate := NewCertificate(string(certificateAsString))
	assert.True(t, certificate.IsValid(firstOfJune2024UTCAt1AM))
}

func TestFindsACertificate(t *testing.T) {
	certificate, err := os.ReadFile("./test_fixtures/RootCA.pem")
	if err != nil {
		t.Fail()
	}

	certificates := GetCertificatesFromString(string(certificate))

	assert.Len(t, certificates, 1)
}

func TestFindsMultipleCertificates(t *testing.T) {
	certificate, err := os.ReadFile("./test_fixtures/RootCA.pem")
	if err != nil {
		t.Fail()
	}

	certificates := certificate
	certificates = append(certificates, certificate...)

	foundCertificates := GetCertificatesFromString(string(certificates))

	assert.Len(t, foundCertificates, 2)
}

func TestFindsNoCertificates(t *testing.T) {
	certificates := GetCertificatesFromString("")

	assert.Empty(t, certificates)
}
