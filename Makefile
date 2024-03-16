.PHONY: run view_build help

run: ## Run the application
	@go run ./cmd/my_authn_authz/main.go

view_build: ## Build the view files
	@cd ./view && yarn build

help: ## Show Makefile options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'