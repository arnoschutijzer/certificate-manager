package main

import (
	"fmt"
	"time"

	"github.com/algleymi/certificate-manager/internal/onepassword"
)

func main() {
	o, err := onepassword.NewOnePassword()

	if err != nil {
		panic(err)
	}

	inAMonth := time.Now().AddDate(0, 1, 0)

	outdatedCertificates, err := o.FindCertificates(inAMonth)

	if err != nil {
		panic(err)
	}

	if len(outdatedCertificates) == 0 {
		fmt.Println("no outdated certificates found!")
		return
	}

	fmt.Println("outdated certificates found!")
	for _, v := range outdatedCertificates {
		fmt.Printf("%s (%s)\n", v.Subject, v.CustomName)
	}
}
