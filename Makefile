.PHONY: build test fmt vet

build:
	mkdir -p bin
	go build -o bin/go-tax-calculator ./cmd/server/main.go

run: build
	./bin/go-tax-calculator

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

