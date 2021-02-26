BIN_DIR               = bin
BIN                   = guestbook
CMD_DIR               = cmd
CMD                   = guestbook
GOLANGCI_LINT         = $(shell pwd)/bin/golangci-lint
GOLANGCI_LINT_VERSION = v1.37.1

ECHO                  = echo "  "
DOCKER                = docker
RM                    = rm
GO                    = go
SHELL                 = /usr/bin/env bash -o pipefail

GOBUILD               = $(GO) build
CGO                   = CGO_ENABLED=0 GOOS=linux
LDFLAGS               = -ldflags="-s -w"
GO_SRC                = $(shell find ./ -name '*.go')
UPPER                 = $(shell echo '$1' | tr '[:lower:]' '[:upper:]')
FORMAT                = sed "s/^/    /g"

ifneq ($(V),1)
	Q = @
endif

IMG ?= guestbook-go:latest

all: help

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; \
		printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: build
build: fmt lint ## Build executable
	$(Q)$(ECHO) "GO" $(call UPPER, $@)
	$(Q)$(CGO) $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BIN) $(CMD_DIR)/$(CMD)/*.go 2>&1 | $(FORMAT)

.PHONY: run
run: fmt lint ## Run against the configured Kubernetes cluster in ~/.kube/config
	$(Q)$(ECHO) "GO" $(call UPPER, $@)
	$(Q)$(GO) run $(CMD_DIR)/$(CMD)/main.go

.PHONY: fmt
fmt: ## Run go fmt against code
	$(Q)$(ECHO) "GO" $(call UPPER, $@)
	$(Q)$(GO) fmt ./... 2>&1 | $(FORMAT)

$(GOLANGCI_LINT):
	$(Q)$(ECHO) "GOLANGCI_LINT"
	$(Q) curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell dirname $(GOLANGCI_LINT)) $(GOLANGCI_LINT_VERSION) 2>&1 | $(FORMAT)

.PHONY: lint
lint: ${GOLANGCI_LINT} ## Run golangci-lint linter
	$(Q)$(ECHO) $(call UPPER, $@)
	$(Q)$(GOLANGCI_LINT) run --color always | $(FORMAT)

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Run golangci-lint linter and perform fixes
	$(Q)$(ECHO) $(call UPPER, $@)
	$(Q)$(GOLANGCI_LINT) run --fix --color always 2>&1 | $(FORMAT)


##@ Publishing
.PHONY: build-docker
build-docker: build ## Build the docker image
	docker build . -t ${IMG}
