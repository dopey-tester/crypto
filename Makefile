# Set V to 1 for verbose output from the Makefile
Q=$(if $V,,@)
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: lint test

ci: test

.PHONY: all ci

#########################################
# Bootstrapping
#########################################

bootstrap:
	$Q go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: bootstrap

#########################################
# Test
#########################################

test:
	$Q $(GOFLAGS) go test -coverprofile=coverage.out ./...

race:
	$Q $(GOFLAGS) go test -race ./...

.PHONY: test race

#########################################
# Linting
#########################################

fmt:
	$Q gofmt -s -l -w $(SRC)

lint:
	$Q LOG_LEVEL=error golangci-lint run --timeout=30m

.PHONY: lint fmt

#########################################
# Go generate
#########################################

generate:
	$Q go generate ./...

.PHONY: generate
