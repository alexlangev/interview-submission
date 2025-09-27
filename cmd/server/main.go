package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexlangev/interview-submission/internal/api"
)

func main() {
	fmt.Println("Smoke")

	const apiURL = "http://localhost:5001"
	const port = ":8080"

	// prov := client.NewClient(apiURL, nil)

	r := api.NewRouter()

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	log.Printf("HTTP server listening on %s (api-url=%s)", port, apiURL)
	log.Fatal(srv.ListenAndServe())
}
