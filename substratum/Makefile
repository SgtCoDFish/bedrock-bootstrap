SHELL := /bin/bash

NAME=substratum
VERSION := $(shell cat VERSION.txt)

GO = go
GOLIST := $(shell $(GO) list -mod=vendor ./... | grep -v vendor)
GOLINT = golint
GOSTATICCHECK = staticcheck
GOLANGCILINT = golangci-lint

.PHONY: ci
ci: clean test fmt lint staticcheck golangci-lint verify-vendor bin/$(NAME)

.PHONY: clean
clean:
	@echo "+ $@"
	@rm -rf bin

.PHONY: test
test:
	@echo "+ $@"
	@$(GO) test -mod=vendor $(GOLIST)

.PHONY: fmt
fmt:
	@echo "+ $@"
	@if [[ ! -z "$(shell gofmt -l -s . | grep -v vendor | tee /dev/stderr)" ]]; then exit 1; fi

.PHONY: staticcheck
staticcheck:
	@echo "+ $@"
	@if [[ ! -z "$(shell $(GOSTATICCHECK) $(GOLIST) | tee /dev/stderr)" ]]; then exit 1; fi

.PHONY: lint
lint:
	@echo "+ $@"
	@if [[ ! -z '$(shell $(GOLINT) ./... | grep -v vendor | tee /dev/stderr)' ]]; then exit 1; fi

.PHONY: golangci-lint
golangci-lint:
	@echo "+ $@"
	@if [[ ! -z '$(shell $(GOLANGCILINT) run | tee /dev/stderr)' ]]; then exit 1; fi

.PHONY: verify-vendor
verify-vendor:
	@echo "+ $@"
	@$(GO) mod verify > /dev/null

.PHONY: build
build: bin/$(NAME)

bin/$(NAME): $(wildcard *.go) $(wildcard cmd/$(NAME)/*.go) $(wildcard autotest/*.go) $(wildcard autotest/uart_rxxd/*.go)
	@mkdir -p bin
	$(GO) build -mod=vendor -o $@ cmd/$(NAME)/*.go

.PHONY: docker-build-env
docker-build-env: VERSION.txt Dockerfile
	@if [[ -z '$(shell docker image ls --quiet substratum-build:$(VERSION))' ]]; then \
		docker build -t substratum-build:$(VERSION) . ;\
	fi

.PHONY: docker-ci
docker-ci: docker-build-env
	mkdir -p bin BUILD
	docker run --rm -v "$$PWD:/work" substratum-build:$(VERSION) make ci
