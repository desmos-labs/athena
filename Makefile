VERSION := $(shell echo $(shell git describe --always) | sed 's/^v//')
TENDERMINT_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')
COMMIT := $(shell git log -1 --format='%H')
BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)

export GO111MODULE = on

all: lint test-unit install

###############################################################################
###                                Build flags                              ###
###############################################################################

# These lines here are essential to include the muslc library for static linking of libraries
# (which is needed for the wasmvm one) available during the build. Without them, the build will fail.
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# Process linker flags
ldflags = -X 'github.com/forbole/juno/v5/cmd.Version=$(VERSION)' \
 	-X 'github.com/forbole/juno/v5/cmd.Commit=$(COMMIT)'

ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

###############################################################################
###                                  Build                                  ###
###############################################################################
BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

.PHONY: build install

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

clean:
	rm -rf $(BUILDDIR)/

.PHONY: go-mod-cache go.sum clean

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

coverage:
	@echo "Viewing test coverage..."
	@go tool cover --html=coverage.out

stop-docker-test:
	@echo "Stopping Docker container..."
	@docker stop djuno-test-db || true && docker rm djuno-test-db || true

start-docker-test: stop-docker-test
	@echo "Starting Docker container..."
	@docker run --name djuno-test-db -e POSTGRES_USER=djuno -e POSTGRES_PASSWORD=password -e POSTGRES_DB=djuno -d -p 6432:5432 postgres

test-unit: start-docker-test
	@echo "Executing unit tests..."
	@go test -mod=readonly -v -coverprofile coverage.txt ./...

.PHONY: coverage stop-docker-test start-docker-test test-unit

###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	@echo "--> Running linter"
	@go run $(golangci_lint_cmd) run --timeout=10m

lint-fix:
	@echo "--> Running linter"
	@go run $(golangci_lint_cmd) run --fix --out-format=tab --issues-exit-code=0

.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' -not -path "./venv" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' -not -path "./venv" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' -not -path "./venv" | xargs goimports -w -local github.com/desmos-labs/djuno
.PHONY: format

.PHONY: lint lint-fix format
