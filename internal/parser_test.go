package internal

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParsesASingleCertificate(t *testing.T) {
	certificate, err := os.ReadFile("./test_fixtures/RootCA.pem")
	if err != nil {
		t.Fail()
	}

	firstOfJune2024UTCAt1AM := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)

	assert.True(t, isValidCertificate(certificate, firstOfJune2024UTCAt1AM))
}
