package main

import "github.com/algleymi/certificate-manager/internal/onepassword"

func main() {
	o, err := onepassword.NewOnePassword()

	if err != nil {
		panic(err)
	}

	o.FindCertificates()

	// inTwoMonths := time.Now().AddDate(0, 2, 0)
	// before := time.Now()
	// certificates, err := onepassword.FindCertificatesOlderThanDate(inTwoMonths)
	// fmt.Printf("Treated all certificates, took %f\n", time.Since(before).Seconds())

	// numberOfCertificates := len(certificates)

	// if numberOfCertificates == 0 {
	// 	fmt.Println("No outdated certificates, nice!")
	// 	return
	// }

	// slices.SortFunc(certificates, func(a, b domain.Certificate) int {
	// 	return cmp.Compare(strings.ToLower(a.CustomName), strings.ToLower(b.CustomName))
	// })

	// fmt.Printf("Found %d outdated certificates\n", numberOfCertificates)
	// for _, v := range certificates {
	// 	fmt.Printf("%s - %s\n", v.CustomName, v.Subject)
	// }
}
