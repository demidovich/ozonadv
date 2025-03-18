.PHONY: vendor

GIT_TAG ?= $(shell git tag)

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

vendor: ## Go mod vendor
	go mod tidy
	go mod vendor

fetch:
	go run ./cmd/ozonadv.go fetch

# ---------------------------------------------------------------------------------------------------------------------

build-win: ## Build windows
	GOOS=windows GOARCH=amd64 go build -o build/cian-street-url-$(GIT_TAG).exe cian_street_url.go

build-mac: ## Build mac
	GOOS=darwin GOARCH=amd64 go build -o build/cian-street-url-mac-$(GIT_TAG) cian_street_url.go
