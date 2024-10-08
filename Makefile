############################################################################################################
# Few rules to follow when editing this file:
#
# 1. Shell commands must be indented with a tab
# 2. Before each target add ".PHONY: target_name" to disable default file target
# 3. Add target description prefixed with "##" on the same line as target definition for "help" target to work
# 4. Be aware that each make command is executed in separate shell
#
# Tips:
#
# * Add an @ sign to suppress output of the command that is executed
# * Define variable like: VAR := value
# * Reference variable like: $(VAR)
# * Reference environment variables like: $(ENV_VARIABLE)
#
#############################################################################################################
.DELETE_ON_ERROR:
# TOOLS
GOLANGCILINT_VERSION = v1.60.1
GOFUMPT_VERSION		 = v0.6.0
GOTESTSUM_VERSION    = v1.12.0
# Shell
.SHELLFLAGS 	:= -eu -o pipefail -c
SHELL			= /bin/bash
# OUT(s)
BIN 			:= .bin
BUILDOUT		= ./bin
# BUILD FLAGS
MAKEFLAGS 		+= --warn-undefined-variables
MAKEFLAGS 		+= --no-builtin-rules
MAKEFLAGS		+= --no-print-directory
GOPATH			?= ${shell go env GOPATH}
export GOBIN 	:= $(abspath $(BIN))
export GOFLAGS 	:= -ldflags=-linkmode=internal
# ARGS etc ...
ARGS 			:=


BOLD = \033[1m
CLEAR = \033[0m
CYAN = \033[36m
GREEN = \033[32m

##@: Default
.PHONY: help
help: ## Display this help
	@awk '\
		BEGIN {FS = ":.*##"; printf "Usage: make $(CYAN)<target>$(CLEAR)\n"} \
		/^[0-9a-zA-Z_\-\/]+?:[^#]*?## .*$$/ { printf "\t$(CYAN)%-20s$(CLEAR) %s\n", $$1, $$2 } \
		/^##@/ { printf "$(BOLD)%s$(CLEAR)\n", substr($$0, 5); }' \
		$(MAKEFILE_LIST)

##@: Build/Run

all: clean go/tidy go/format test test_e2e lint release debug  ## ALL, builds the world

.PHONY: clean
clean: ## Clean-up build artifacts
	@echo -e "$(CYAN)--- clean...$(CLEAR)"
	@go clean
	@rm -rf ${BUILDOUT}

.PHONY: test
test: clean go/test ## Runs all tests

.PHONY: test_e2e
test_e2e: clean go/test_e2e ## Runs all e2e tests

.PHONY: bench
bench: clean go/benchmark ## Runs all benchmarks

.PHONY: lint
lint: go/lint ## Runs all linters

.PHONY: release
release: go/release ## Build RELEASE (debug off)

.PHONY: debug
debug: go/debug ## Build DEBUG (debug on)

.PHONY: run
run: ## Runs golox-vm. Use 'make ARGS="script.lox" run' to pass arguments
	@echo -e "$(CYAN)--- run ...$(CLEAR)"
	go run ./main.go $(ARGS)

.PHONY: rund
rund: ## Runs goloxd-vm. Use 'make ARGS="script.lox" run' to pass arguments
	@echo -e "$(CYAN)--- run ...$(CLEAR)"
	go run -race -tags debug ./main.go $(ARGS)

###@: Go
.PHONY: go/format
go/format: $(BIN)/gofumpt ### Format all go files
	@echo -e "$(CYAN)--- format go files...$(CLEAR)"
	$(BIN)/gofumpt -w .

go/tidy: go.mod go.sum ### Tidy all Go dependencies
	@echo -e "$(CYAN)--- tidy go dependencies...$(CLEAR)"
	go mod tidy -v -x

.PHONY: go/lint
go/lint: $(BIN)/golangci-lint ### Lints the codebase using golangci-lint
	@echo -e "$(CYAN)--- lint codebase...$(CLEAR)"
	$(BIN)/golangci-lint run --modules-download-mode=readonly --config .golangci.yml

.PHONY: go/test
go/test: $(BIN)/gotestsum ### Runs unit tests
	@echo -e "$(CYAN)--- go test ...$(CLEAR)"
	@$(BIN)/gotestsum --debug --format-hide-empty-pkg --format=testdox -- -shuffle=on -race -timeout=60s -count 1 -parallel 3 -v ./internal/...

.PHONY: go/test_e2e
go/test_e2e: $(BIN)/gotestsum ### Runs e2e tests
	@echo -e "$(CYAN)--- go test_e2e ...$(CLEAR)"
	@$(BIN)/gotestsum --debug --format-hide-empty-pkg --format=testdox -- -shuffle=on -race -timeout=60s -count 1 -parallel 3 -v ./test_e2e/...

.PHONY: go/benchmark
go/benchmark: ### Runs unit tests
	@echo -e "$(CYAN)--- go benchmark ...$(CLEAR)"
	go test -v -bench=. -run='^#' -benchtime=5s -timeout=30m ./test_e2e/...

.PHONY: go/release
go/release: ### Build
	@echo -e "$(CYAN)--- go/build ...$(CLEAR)"
	go build -tags release -o ${BUILDOUT}/golox-vm ./main.go

.PHONY: go/debug
go/debug: ### Build
	@echo -e "$(CYAN)--- go/build ...$(CLEAR)"
	go build -tags debug -o ${BUILDOUT}/goloxd-vm ./main.go

# TOOLS
$(BIN)/golangci-lint: Makefile
	@mkdir -p $(@D)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

$(BIN)/gofumpt: Makefile
	@mkdir -p $(@D)
	go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

$(BIN)/gotestsum: Makefile
	@mkdir -p $(@D)
	go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)
