GOCMD=go
BINARY_DIR=.
BUILD_PKG=github.com/will-steinleitner/will-steinleitner.de/backend/internal/common/build
VERSION ?= $(if $(APP_VERSION),$(APP_VERSION),$(shell git describe --tags --abbrev=0 2>/dev/null || echo development))
GIT_HASH ?= $(if $(APP_GIT_HASH),$(APP_GIT_HASH),$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown))
BUILD_LDFLAGS=-ldflags "-X $(BUILD_PKG).Version=$(VERSION) -X $(BUILD_PKG).GitHash=$(GIT_HASH)"

.PHONY: deps
deps: ## Install all Go dependencies
	$(GOCMD) mod tidy

.PHONY: clean
clean: ## Clean up all build artifacts
	$(GOCMD) clean
	rm -f $(BINARY_DIR)/server $(BINARY_DIR)/server.exe

.PHONY: build
build: build-server ## Build all

.PHONY: build-server
build-server: ## Build server application
	mkdir -p $(BINARY_DIR)
	$(GOCMD) build $(BUILD_LDFLAGS) -o $(BINARY_DIR)/server ./cmd/server/

.PHONY: deploy
deploy: build-server ## Build and deploy to Hetzner
	ssh hetzner "mkdir -p /root/projects/will-steinleitner.de"
	scp $(BINARY_DIR)/server hetzner:/root/projects/will-steinleitner.de/server

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
