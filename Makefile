.PHONY: dependencies_start dependencies_stop migrate_up migrate_down migrate_reset insert_data clear_data run view_build development_start try_auth_flow help

dependencies_start: ## Start the postgres and redis
	@docker compose --env-file .env.local -f ./docker/docker-compose.local.yaml up -d

dependencies_stop: ## Stop the postgres and redis
	@docker compose --env-file .env.local -f ./docker/docker-compose.local.yaml down

migrate_up: ## Run the migration
	@bash -c 'source .env.local && migrate -path ./db/migration -database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@$${POSTGRES_HOST}:$${POSTGRES_PORT}/$${POSTGRES_DB}?sslmode=$${POSTGRES_SSL_MODE}" up'

migrate_down: ## Rollback the migration
	@bash -c 'source .env.local && migrate -path ./db/migration -database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@$${POSTGRES_HOST}:$${POSTGRES_PORT}/$${POSTGRES_DB}?sslmode=$${POSTGRES_SSL_MODE}" down'

migrate_reset: ## Reset the migration
	$(MAKE) migrate_down
	$(MAKE) migrate_up

insert_data: ## Insert data into the database
	@bash -c 'source .env.local && PGPASSWORD=$${POSTGRES_PASSWORD} psql -h $${POSTGRES_HOST} -U $${POSTGRES_USER} -d $${POSTGRES_DB} -f ./db/data/insert.sql'

clear_data: ## Clear data from the database
	@bash -c 'source .env.local && PGPASSWORD=$${POSTGRES_PASSWORD} psql -h $${POSTGRES_HOST} -U $${POSTGRES_USER} -d $${POSTGRES_DB} -f ./db/data/clear.sql'

run: ## Run the application
	@ECHO_MODE=local go run ./cmd/lgtmeme/main.go

clean: ## Clean the binary
	go clean -cache -modcache

view_build: ## Build the view files
	@cd ./view && yarn install && yarn build

development_start: ## Start the development environment
	@bash -c './script/start_development.sh'

try_auth_flow: ## Try the authentication and authorization flow
	@bash -c './script/try_auth_flow.sh'

migrate_up_for_prod: ## Run the migration for production
	@bash -c 'source .env.prod && migrate -path ./db/migration -database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@$${POSTGRES_HOST}:$${POSTGRES_PORT}/$${POSTGRES_DB}?sslmode=$${POSTGRES_SSL_MODE}" up'

migrate_down_for_prod: ## Rollback the migration for production
	@bash -c 'source .env.prod && migrate -path ./db/migration -database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@$${POSTGRES_HOST}:$${POSTGRES_PORT}/$${POSTGRES_DB}?sslmode=$${POSTGRES_SSL_MODE}" down'

help: ## Show Makefile options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'