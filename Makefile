APP_NAME = ozonadv
GIT_TAG ?= $(shell git tag)

.PHONY: vendor clean

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

vendor: ## Go mod vendor
	go mod tidy
	sleep 1
	go mod vendor

lint: ## Run static tests
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run ./...

test: ## Run tests
	go test ./...

# ---------------------------------------------------------------------------------------------------------------------

build-linux: ## Build Linux
	GOOS=linux GOARCH=amd64 go build -o build/${APP_NAME}_Linux main.go

build-darwin: ## Build Mac
	GOOS=darwin GOARCH=amd64 go build -o build/${APP_NAME}_Darwin main.go

build-windows: ## Build Windows
	GOOS=windows GOARCH=amd64 go build -o build/${APP_NAME}_Windows.exe main.go

build-ci-linux: ## Build CI Mac
	GOOS=darwin GOARCH=amd64 go build -o build/${APP_NAME}_${GIT_TAG}_Linux_x86_64 main.go

build-ci-darwin: ## Build CI Mac
	GOOS=darwin GOARCH=amd64 go build -o build/${APP_NAME}_${GIT_TAG}_Darwin_x86_64 main.go

build-ci-windows: ## Build CI Windows
	GOOS=windows GOARCH=amd64 go build -o build/${APP_NAME}_${GIT_TAG}_Windows_x86_64}.exe main.go

clean: ## Clean build directory.
	rm -f ./build/*
