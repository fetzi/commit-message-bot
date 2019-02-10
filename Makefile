PROJECT := commit-message-bot
BINDIR  := $(CURDIR)/bin
VERSION := $(shell cat VERSION)
XC_OS 	:= linux darwin
XC_ARCH := 386 amd64 arm

export XC_OS
export XC_ARCH
export VERSION

.PHONY: all
all: help


.PHONY: help
help:
	@echo "make test - run go test"
	@echo "make build - build $(PROJECT) for follwing OS-ARCH constilations: $(XC_OS) / $(XC_ARCH) "
	@echo "make build-dev - build $(PROJECT) for OS-ARCH set by GOOS and GOARCH env variables"
	@echo "make dist - build and create packages with hashsums"


.PHONY: build
build:
	@scripts/build.sh

.PHONY: build-dev
build-dev:
	@scripts/build.sh dev

.PHONY: dist
dist:
	@scripts/dist.sh