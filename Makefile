# This file is modified from the Makefile of
# bsky written in Go by Yasuhiro Matsumoto (a.k.a. mattn)
# <https://github.com/mattn/bsky/>
# Permission granted for reuse by the MIT License of bsky.
# Kenji Rikitake distributes this file
# under the MIT License of gossn.

BIN := gossn
VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS := "-s -w -X main.revision=$(CURRENT_REVISION)"
GOOS := $(shell go env GOOS)
GOBIN ?= $(shell go env GOPATH)/bin
export GO111MODULE=on

.PHONY: all
all: clean build

.PHONY: build
build:
	go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) .

.PHONY: release
release:
	go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) .
	zip -r gossn-$(GOOS)-$(VERSION).zip $(BIN)

.PHONY: install
install:
	go install -ldflags=$(BUILD_LDFLAGS) .

.PHONY: show-version
show-version: $(GOBIN)/gobump
	gobump show -r .

$(GOBIN)/gobump:
	go install github.com/x-motemen/gobump/cmd/gobump@latest

.PHONY: test
test: build
	go test -v ./...

.PHONY: clean
clean:
	go clean

.PHONY: bump
bump: $(GOBIN)/gobump
ifneq ($(shell git status --porcelain),)
	$(error git workspace is dirty)
endif
ifneq ($(shell git rev-parse --abbrev-ref HEAD),main)
	$(error current branch is not main)
endif
	@gobump up -w .
	git commit -am "bump up version to $(VERSION)"
	git tag "v$(VERSION)"
	git push origin main
	git push origin "refs/tags/v$(VERSION)"
