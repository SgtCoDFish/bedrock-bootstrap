SHELL := /bin/bash

NAME=substratum
SS_ASM_NAME=ss-asm

VERSION := $(shell cat VERSION.txt)

GO = go
GOLIST := $(shell $(GO) list -mod=vendor ./... | grep -v vendor)
GOLINT = golint

.PHONY: ci
ci: clean test fmt lint golangci-lint verify-vendor bin/$(NAME)

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

.PHONY: lint
lint:
	@echo "+ $@"
	@if [[ ! -z '$(shell $(GOLINT) ./... | grep -v vendor | tee /dev/stderr)' ]]; then exit 1; fi

.PHONY: golangci-lint
golangci-lint: devel/bin/golangci-lint
	@echo "+ $@"
	$< run --modules-download-mode vendor

.PHONY: verify-vendor
verify-vendor:
	@echo "+ $@"
	@$(GO) mod verify > /dev/null

.PHONY: build
build: bin/$(NAME) bin/$(SS_ASM_NAME)

bin:
	@mkdir -p bin

bin/$(NAME): $(wildcard *.go) $(wildcard cmd/$(NAME)/*.go) $(wildcard cmd/ss-asm/*.go) $(wildcard autotest/*.go) $(wildcard autotest/uart_rxxd/*.go) | bin
	$(GO) build -mod=vendor -o $@ cmd/$(NAME)/*.go

bin/$(SS_ASM_NAME): $(wildcard *.go) $(wildcard cmd/$(NAME)/*.go) $(wildcard cmd/ss-asm/*.go) $(wildcard cmd/ss-asm/main/*.go) | bin
	$(GO) build -mod=vendor -o $@ cmd/$(SS_ASM_NAME)/main/*.go

.PHONY: docker-build-env
docker-build-env: VERSION.txt Dockerfile
	@if [[ -z '$(shell docker image ls --quiet substratum-build:$(VERSION))' ]]; then \
		docker build -t substratum-build:$(VERSION) . ;\
	fi

.PHONY: docker-ci
docker-ci: docker-build-env
	mkdir -p bin BUILD
	docker run --rm -v "$$PWD:/work" substratum-build:$(VERSION) make ci

.PHONY: install
install:
	sudo cp bin/$(NAME) /usr/bin/$(NAME)

devel/bin/golangci-lint:
	@mkdir -p devel/bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./devel/bin v1.41.1

