PACKAGES ?= $(shell go list ./... | grep -Fv -e mocks )

gobuild: ## Build Orchestrate Go binary
	@GOOS=linux GOARCH=amd64 go build -i -o ./build/bin/key-manager

run-coverage:
	@sh scripts/coverage.sh $(PACKAGES)

coverage: run-coverage
	@$(OPEN) build/coverage/coverage.html 2>/dev/null