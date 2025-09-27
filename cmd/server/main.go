package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexlangev/interview-submission/client"
)

func main() {
	fmt.Println("Smoke")

	BaseURL := "http://localhost:5001"
	httpClient := &http.Client{}

	client := &client.Client{
		BaseURL:    BaseURL,
		HTTPClient: httpClient,
	}

	brackets, err := client.GetTaxBrackets(2022)
	if err != nil {
		log.Fatalf("GetTaxBrackets: %v", err)
	}

	fmt.Println(brackets)

}
