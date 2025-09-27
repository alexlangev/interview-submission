package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexlangev/interview-submission/client"
)

func main() {
	fmt.Println("Smoke")

	baseURL := "http://localhost:5001"
	httpClient := &http.Client{Timeout: 5 * time.Second}
	c := client.NewClient(baseURL, httpClient)

	brackets, err := c.GetTaxBrackets(2022)
	if err != nil {
		log.Fatalf("GetTaxBrackets: %v", err)
	}

	out, _ := json.MarshalIndent(map[string]any{"tax_brackets": brackets}, "", "  ")
	fmt.Println(string(out))
}
