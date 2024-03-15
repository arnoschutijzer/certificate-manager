package internal

import (
	"crypto/x509"
	"encoding/pem"
	"strings"
	"time"
)

func getCertificatesFromString(secret string) []string {
	certificates := []string{}
	remainingSubstrings := secret

	startCertificateTemplate := "-----BEGIN CERTIFICATE-----"
	endCertificateTemplate := "-----END CERTIFICATE-----"

	for true {
		indexOfBegin := strings.Index(remainingSubstrings, startCertificateTemplate)
		indexOfEnd := strings.Index(remainingSubstrings, endCertificateTemplate)

		if indexOfBegin == -1 {
			return certificates
		}

		certificate := remainingSubstrings[indexOfBegin : indexOfEnd+len(endCertificateTemplate)]

		certificates = append(certificates, certificate)
		remainingSubstrings = remainingSubstrings[indexOfEnd+len(endCertificateTemplate):]
	}

	return certificates
}

func isValidCertificate(certificate []byte, aDate time.Time) bool {
	block, _ := pem.Decode(certificate)
	cert, _ := x509.ParseCertificate(block.Bytes)

	notAfter := cert.NotAfter
	notBefore := cert.NotBefore

	return notBefore.Before(aDate) && notAfter.After(aDate)
}
