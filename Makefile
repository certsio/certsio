APP_NAME := certsio
UNAME := $(shell uname)
PWD := $(shell pwd)
CMD_DIR := cmd/certsio

build:
	@if [ "$(UNAME)" = "Darwin" ]; then \
		make mac; \
	else \
		make linux; \
	fi

mac:
	@go build -o $(APP_NAME) $(PWD)/$(CMD_DIR)

linux:
	@echo "Building $(APP_NAME) for Linux"
	@GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) $(PWD)/$(CMD_DIR)

test: fmt
	@go test -v ./...

fmt:
	@gofumpt -l -w .

lint:
	@golangci-lint run ./...

.PHONY: build fmt lint test
