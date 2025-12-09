BINARY_NAME=payjp
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-ldflags "-X github.com/payjp/payjp-cli/cmd.Version=${VERSION}"

.PHONY: all build clean test install lint fmt deps help

all: build

## build: Build the binary
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} .

## install: Install the binary
install:
	go install ${LDFLAGS} .

## test: Run tests
test:
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## clean: Remove build artifacts
clean:
	rm -f ${BINARY_NAME}
	rm -f coverage.out coverage.html
	rm -rf dist/

## lint: Run linter
lint:
	golangci-lint run ./...

## fmt: Format code
fmt:
	go fmt ./...
	goimports -w .

## deps: Download dependencies
deps:
	go mod download
	go mod tidy

## build-all: Build for all platforms
build-all: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe .

## release: Create release archives
release: build-all
	cd dist && tar -czf ${BINARY_NAME}-linux-amd64.tar.gz ${BINARY_NAME}-linux-amd64
	cd dist && tar -czf ${BINARY_NAME}-linux-arm64.tar.gz ${BINARY_NAME}-linux-arm64
	cd dist && tar -czf ${BINARY_NAME}-darwin-amd64.tar.gz ${BINARY_NAME}-darwin-amd64
	cd dist && tar -czf ${BINARY_NAME}-darwin-arm64.tar.gz ${BINARY_NAME}-darwin-arm64
	cd dist && zip ${BINARY_NAME}-windows-amd64.zip ${BINARY_NAME}-windows-amd64.exe

## version: Show version
version:
	@echo ${VERSION}

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' ${MAKEFILE_LIST} | column -t -s ':' | sed 's/^/  /'
