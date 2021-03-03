gobuild: ## Build Orchestrate Go binary
	@GOOS=linux GOARCH=amd64 go build -i -o ./build/bin/key-manager
