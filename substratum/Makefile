MAKEFLAGS += --warn-undefined-variables --no-builtin-rules
SHELL := /usr/bin/env bash
.SHELLFLAGS := -uo pipefail -c
.DEFAULT_GOAL := ci
.DELETE_ON_ERROR:
.SUFFIXES:

BINDIR := _bin

NAME := substratum
SS_ASM_NAME := ss-asm
SS_AUTOTEST_NAME := ss-autotest

VERSION := $(shell cat VERSION.txt)

GOLANGCILINT_VERSION := v1.55.2

GO := go
GOLIST := $(shell $(GO) list -mod=vendor ./... | grep -v vendor)

.PHONY: ci
ci: test fmt golangci-lint verify-vendor build

.PHONY: clean
clean:
	@rm -rf $(BINDIR)

.PHONY: test
test:
	@$(GO) test -mod=vendor $(GOLIST)

.PHONY: fmt
fmt:
	@echo "+ $@"
	@if [[ ! -z "$(shell gofmt -l -s . | grep -v vendor | tee /dev/stderr)" ]]; then exit 1; fi

.PHONY: golangci-lint
golangci-lint: $(BINDIR)/golangci-lint-$(GOLANGCILINT_VERSION)/golangci-lint
	$< run --modules-download-mode vendor

.PHONY: verify-vendor
verify-vendor:
	$(GO) mod verify > /dev/null

.PHONY: build
build: $(BINDIR)/$(NAME) $(BINDIR)/$(SS_ASM_NAME) $(BINDIR)/$(SS_AUTOTEST_NAME)

$(BINDIR)/$(NAME): $(wildcard *.go) $(wildcard cmd/$(NAME)/*.go) $(wildcard cmd/$(SS_ASM_NAME)/*.go) $(wildcard cmd/$(SS_AUTOTEST_NAME)/*.go) $(wildcard autotest/*.go) | $(BINDIR)
	$(GO) build -mod=vendor -o $@ cmd/$(NAME)/*.go

$(BINDIR)/$(SS_ASM_NAME): $(wildcard *.go) $(wildcard cmd/$(SS_ASM_NAME)/*.go) $(wildcard cmd/$(SS_ASM_NAME)/main/*.go) | $(BINDIR)
	$(GO) build -mod=vendor -o $@ cmd/$(SS_ASM_NAME)/main/*.go

$(BINDIR)/$(SS_AUTOTEST_NAME): $(wildcard *.go) $(wildcard cmd/$(SS_AUTOTEST_NAME)/*.go) $(wildcard cmd/$(SS_AUTOTEST_NAME)/main/*.go) | $(BINDIR)
	$(GO) build -mod=vendor -o $@ cmd/$(SS_AUTOTEST_NAME)/main/*.go

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
	sudo cp $(BINDIR)/$(NAME) /usr/bin/$(NAME)

$(BINDIR)/golangci-lint-$(GOLANGCILINT_VERSION)/golangci-lint:
	@mkdir -p $(@D)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(@D)/ $(GOLANGCILINT_VERSION)

$(BINDIR):
	@mkdir -p $(BINDIR)

