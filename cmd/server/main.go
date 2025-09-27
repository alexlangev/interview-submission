package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/alexlangev/interview-submission/internal/api"
	"github.com/alexlangev/interview-submission/internal/client"
	"github.com/alexlangev/interview-submission/internal/core"
)

func main() {
	apiURL := flag.String("api-url", "http://localhost:5001", "tax API base URL")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	prov := client.NewClient(*apiURL, nil)
	calc := core.NewCalculator(prov)

	r := api.NewRouter(calc)

	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	log.Printf("HTTP server listening on %s", *addr)
	log.Fatal(srv.ListenAndServe())
}
