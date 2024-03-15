package internal

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

func getCertificatesFromString(secret string) {

}

func isValidCertificate(certificate []byte, aDate time.Time) bool {
	block, _ := pem.Decode(certificate)
	cert, _ := x509.ParseCertificate(block.Bytes)
	fmt.Println(cert.Subject.Organization)
	notAfter := cert.NotAfter
	notBefore := cert.NotBefore

	return notBefore.Before(aDate) && notAfter.After(aDate)
}
