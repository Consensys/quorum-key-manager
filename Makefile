GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*" | egrep -v "^\./\.go" | grep -v _test.go)
DEPS_HASHICORP = hashicorp hashicorp-init hashicorp-agent
PACKAGES ?= $(shell go list ./... | egrep -v "integration-tests|mocks" )
KEY_MANAGER_SERVICES = key-manager

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OPEN = xdg-open
endif
ifeq ($(UNAME_S),Darwin)
	OPEN = open
endif

.PHONY: all lint integration-tests

lint: ## Run linter to fix issues
	@misspell -w $(GOFILES)
	@golangci-lint run --fix

lint-ci: ## Check linting
	@misspell -error $(GOFILES)
	@golangci-lint run

lint-tools: ## Install linting tools
	@GO111MODULE=on go get github.com/client9/misspell/cmd/misspell@v0.3.4
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.27.0

hashicorp:
	@docker-compose -f deps/docker-compose.yml up --build -d $(DEPS_VAULT)

hashicorp-down:
	@docker-compose -f deps/docker-compose.yml down $(DEPS_VAULT)

deps: hashicorp

down-deps: hashicorp-down

run-acceptance:
	@go test -v -tags acceptance ./acceptance-tests

gobuild:
	@GOOS=linux GOARCH=amd64 go build -i -o ./build/bin/key-manager

run-coverage:
	@sh scripts/coverage.sh $(PACKAGES)

coverage: run-coverage
	@$(OPEN) build/coverage/coverage.html 2>/dev/null

dev: deps gobuild
	@docker-compose -f ./docker-compose.yml up --build -d $(KEY_MANAGER_SERVICES)
	
down-dev: down-deps
	@docker-compose -f ./docker-compose.yml down $(KEY_MANAGER_SERVICES)

run: gobuild
	@build/bin/key-manager run

go-quorum:
	@docker-compose -f deps/go-quorum/docker-compose.yml up -d

stop-go-quorum:
	@docker-compose -f deps/go-quorum/docker-compose.yml stop

down-go-quorum:
	@docker-compose -f deps/go-quorum/docker-compose.yml down --volumes --timeout 0

besu:
	@docker-compose -f deps/besu/docker-compose.yml up -d

stop-besu:
	@docker-compose -f deps/besu/docker-compose.yml stop

down-besu:
	@docker-compose -f deps/besu/docker-compose.yml down --volumes --timeout 0