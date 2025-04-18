.PHONY: vendor

GIT_TAG ?= $(shell git tag)

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

vendor: ## Go mod vendor
	go mod tidy
	sleep 3
	go mod vendor

# ---------------------------------------------------------------------------------------------------------------------

build-linux: ## Build Linux
	GOOS=linux GOARCH=amd64 go build -o build/ozonadv-linux main.go

build-mac: ## Build Mac
	GOOS=darwin GOARCH=amd64 go build -o build/ozonadv-mac main.go

build-win: ## Build Windows
	GOOS=windows GOARCH=amd64 go build -o build/ozonadv.exe main.go

build-ci-linux: ## Build CI Mac
	GOOS=darwin GOARCH=amd64 go build -o build/ozonadv-linux-$(GIT_TAG) main.go

build-ci-mac: ## Build CI Mac
	GOOS=darwin GOARCH=amd64 go build -o build/ozonadv-mac-$(GIT_TAG) main.go

build-ci-win: ## Build CI Windows
	GOOS=windows GOARCH=amd64 go build -o build/ozonadv-win-$(GIT_TAG).exe main.go

