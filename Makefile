VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')
DOCKER := $(shell which docker)

export GO111MODULE = on

all: ci-lint ci-test install

###############################################################################
# Build / Install
###############################################################################

LD_FLAGS = -X github.com/desmos-labs/juno/version.Version=$(VERSION) \
	-X github.com/desmos-labs/juno/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building djuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/djuno.exe ./cmd/djuno
else
	@echo "building djuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/djuno ./cmd/djuno
endif

install: go.sum
	@echo "installing djuno binary..."
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/djuno

###############################################################################
# Tests / CI
###############################################################################

coverage:
	@echo "Viewing test coverage..."
	@go tool cover --html=coverage.out

stop-docker-test:
	@echo "Stopping Docker container..."
	@docker stop djuno-test-db || true && docker rm djuno-test-db || true

start-docker-test: stop-docker-test
	@echo "Starting Docker container..."
	@docker run --name djuno-test-db -e POSTGRES_USER=djuno -e POSTGRES_PASSWORD=password -e POSTGRES_DB=djuno -d -p 5433:5432 postgres

test-unit: start-docker-test
	@echo "Executing unit tests..."
	@go test -mod=readonly -v -coverprofile coverage.txt ./...

lint:
	golangci-lint run --out-format=tab

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0
.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "*.git*" | xargs goimports -w -local github.com/desmos-labs/djuno
.PHONY: format

clean:
	rm -f tools-stamp ./build/**

.PHONY: install build ci-test ci-lint coverage clean start-docker-test
