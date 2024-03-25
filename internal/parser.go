package internal

import (
	"strings"

	"github.com/algleymi/certificate-manager/internal/domain"
)

const (
	START_CERTIFICATE_TEMPLATE = "-----BEGIN CERTIFICATE-----"
	END_CERTIFICATE_TEMPLATE   = "-----END CERTIFICATE-----"
)

func GetCertificatesFromString(secret string, vaultName string) []domain.Certificate {
	// Some secrets are oddly escaped due to them being
	// copy-pasted to Intellij. These are saved with \n directly in the string.
	secret = strings.ReplaceAll(secret, "\\n", "\n")

	certificates := []domain.Certificate{}
	remainingSubstrings := secret

	for true {
		indexOfBegin := strings.Index(remainingSubstrings, START_CERTIFICATE_TEMPLATE)
		indexOfEnd := strings.Index(remainingSubstrings, END_CERTIFICATE_TEMPLATE)

		if indexOfBegin == -1 {
			return certificates
		}

		certificate := remainingSubstrings[indexOfBegin : indexOfEnd+len(END_CERTIFICATE_TEMPLATE)]

		certificates = append(certificates, domain.NewCertificate(certificate, vaultName))
		remainingSubstrings = remainingSubstrings[indexOfEnd+len(END_CERTIFICATE_TEMPLATE):]
	}

	return certificates
}
