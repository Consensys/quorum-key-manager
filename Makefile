GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*" | egrep -v "^\./\.go" | grep -v _test.go)
DEPS_HASHICORP = hashicorp hashicorp-init hashicorp-agent
PACKAGES ?= $(shell go list ./... | egrep -v "acceptance-tests|e2e|mocks|mock" )
KEY_MANAGER_SERVICES = key-manager

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OPEN = xdg-open
endif
ifeq ($(UNAME_S),Darwin)
	OPEN = open
endif

ifneq (,$(wildcard ./.env))
    include .env
    export
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
	@docker-compose -f deps/hashicorp/docker-compose.yml up --build -d $(DEPS_HASHICORP)
	@sleep 2 # Sleep couple seconds to wait token to be created

hashicorp-down:
	@docker-compose -f deps/hashicorp/docker-compose.yml down --volumes --timeout 0

networks:
	@docker network create --driver=bridge hashicorp || true
	@docker network create --driver=bridge --subnet=172.16.237.0/24 besu || true
	@docker network create --driver=bridge --subnet=172.16.238.0/24 quorum || true

down-networks:
	@docker network rm quorum || true
	@docker network rm hashicorp || true

deps: networks hashicorp

down-deps: hashicorp-down

run-acceptance:
	@go test -v -tags acceptance -count=1 ./acceptance-tests

run-e2e:
	@go test -v -tags e2e -count=1 ./e2e

gobuild:
	@GOOS=linux GOARCH=amd64 go build -i -o ./build/bin/key-manager

gobuild-dbg:
	CGO_ENABLED=1 go build -gcflags=all="-N -l" -i -o ./build/bin/key-manager

run-coverage:
	@sh scripts/coverage.sh $(PACKAGES)

coverage: run-coverage
	@$(OPEN) build/coverage/coverage.html 2>/dev/null

dev: gobuild
	@docker-compose -f ./docker-compose.yml up --build $(KEY_MANAGER_SERVICES)	

up: deps gobuild
	@docker-compose -f ./docker-compose.yml up --build -d $(KEY_MANAGER_SERVICES)
	
down:
	@docker-compose -f ./docker-compose.yml down --volumes --timeout 0
	@make down-deps

down-dev:
	@docker-compose -f ./docker-compose.yml down --volumes --timeout 0

run: gobuild
	@build/bin/key-manager run

run-dbg: gobuild-dbg
	@dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./build/bin/key-manager run

go-quorum: networks
	@docker-compose -f deps/go-quorum/docker-compose.yml up -d

stop-go-quorum:
	@docker-compose -f deps/go-quorum/docker-compose.yml stop

down-go-quorum:
	@docker-compose -f deps/go-quorum/docker-compose.yml down --volumes --timeout 0

besu: networks
	@docker-compose -f deps/besu/docker-compose.yml up -d

stop-besu:
	@docker-compose -f deps/besu/docker-compose.yml stop

down-besu:
	@docker-compose -f deps/besu/docker-compose.yml down --volumes --timeout 0
