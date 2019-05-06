SHELL := /bin/bash

NAME=substratum
VERSION := $(shell cat VERSION.txt)

GO = go
GOLIST := $(shell $(GO) list ./... | grep -v vendor)

.PHONY: ci
ci: clean test vet fmt lint staticcheck verify-vendor bin/$(NAME)

.PHONY: clean
clean:
	@echo "+ $@"
	@rm -rf bin

.PHONY: test
test:
	@echo "+ $@"
	@$(GO) test $(GOLIST)

.PHONY: vet
vet:
	@echo "+ $@"
	@if [[ ! -z "$(shell $(GO) vet $(GOLIST))" ]]; then exit 1; fi

.PHONY: fmt
fmt:
	@echo "+ $@"
	@if [[ ! -z "$(shell gofmt -l -s . | grep -v vendor | tee /dev/stderr)" ]]; then exit 1; fi

.PHONY: staticcheck
staticcheck:
	@echo "+ $@"
	@if [[ ! -z "$(shell staticcheck $(GOLIST) | tee /dev/stderr)" ]]; then exit 1; fi


.PHONY: lint
lint:
	@echo "+ $@"
	@if [[ ! -z '$(shell golint ./... | grep -v vendor | tee /dev/stderr)' ]]; then exit 1; fi


.PHONY: verify-vendor
verify-vendor:
	@$(GO) mod verify

.PHONY: build
build: bin/$(NAME)

bin/$(NAME): $(wildcard *.go) $(wildcard cmd/$(NAME)/*.go)
	mkdir -p bin
	$(GO) build -mod=vendor -o $@ cmd/$(NAME)/main.go