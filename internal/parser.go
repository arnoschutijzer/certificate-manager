package internal

import (
	"strings"
)

const (
	START_CERTIFICATE_TEMPLATE = "-----BEGIN CERTIFICATE-----"
	END_CERTIFICATE_TEMPLATE   = "-----END CERTIFICATE-----"
)

func GetCertificatesFromString(secret string, vaultName string) []Certificate {
	// Some secrets are oddly escaped due to them being
	// copy-pasted to Intellij. These are saved with \n directly in the string.
	secret = strings.ReplaceAll(secret, "\\n", "\n")

	certificates := []Certificate{}
	remainingSubstrings := secret

	for true {
		indexOfBegin := strings.Index(remainingSubstrings, START_CERTIFICATE_TEMPLATE)
		indexOfEnd := strings.Index(remainingSubstrings, END_CERTIFICATE_TEMPLATE)

		if indexOfBegin == -1 {
			return certificates
		}

		certificate := remainingSubstrings[indexOfBegin : indexOfEnd+len(END_CERTIFICATE_TEMPLATE)]

		certificates = append(certificates, NewCertificate(certificate, vaultName))
		remainingSubstrings = remainingSubstrings[indexOfEnd+len(END_CERTIFICATE_TEMPLATE):]
	}

	return certificates
}

func DoesSecretContainAnyCertificate(certificate string) bool {
	return strings.Index(certificate, START_CERTIFICATE_TEMPLATE) > -1
}
