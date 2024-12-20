package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func hello() {
	fmt.Println("hello")
}

func main() {
	lambda.Start(hello)
}
