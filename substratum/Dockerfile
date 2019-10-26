FROM golang:1.13-buster

ENV GO111MODULE=on

WORKDIR /work

RUN go get -u golang.org/x/lint/golint && \
	cd /tmp && \
	curl -sSL -o /tmp/golangci-lint.tar.gz https://github.com/golangci/golangci-lint/releases/download/v1.21.0/golangci-lint-1.21.0-linux-amd64.tar.gz && \
	tar xf golangci-lint.tar.gz && \
	mv golangci-lint-1.21.0-linux-amd64/golangci-lint /usr/bin/golangci-lint && \
	chmod +x /usr/bin/golangci-lint && \
	curl -sSL -o /tmp/staticcheck.tar.gz https://github.com/dominikh/go-tools/releases/download/2019.2.3/staticcheck_linux_amd64.tar.gz && \
	tar xf staticcheck.tar.gz && \
	mv staticcheck/staticcheck /usr/bin/staticcheck && \
	chmod +x /usr/bin/staticcheck
