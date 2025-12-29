default: fmt lint install generate

GOCACHE ?= $(CURDIR)/.gocache
export GOCACHE

GOLANGCI_LINT_CACHE ?= $(CURDIR)/.golangci-lint-cache
export GOLANGCI_LINT_CACHE

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
