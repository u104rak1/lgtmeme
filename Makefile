.PHONY: dependencies_start dependencies_stop migrate_up migrate_down migrate_reset insert_data clear_data run view_build help

dependencies_start: ## Start the postgres and redis
	@docker compose --env-file .env.local -f ./docker/docker-compose.local.yaml up -d

dependencies_stop: ## Stop the postgres and redis
	@docker compose --env-file .env.local -f ./docker/docker-compose.local.yaml down

migrate_up: ## Run the migrations
	@bash -c 'source .env.local && migrate -path ./db/migrations -database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@$${POSTGRES_HOST}:$${POSTGRES_PORT}/$${POSTGRES_DB}?sslmode=disable" up'

migrate_down: ## Rollback the migrations
	@bash -c 'source .env.local && migrate -path ./db/migrations -database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@$${POSTGRES_HOST}:$${POSTGRES_PORT}/$${POSTGRES_DB}?sslmode=disable" down'

migrate_reset: ## Reset the migrations
	$(MAKE) migrate_down
	$(MAKE) migrate_up

insert_data: ## Insert data into the database
	@bash -c 'source .env.local && PGPASSWORD=$${POSTGRES_PASSWORD} psql -h $${POSTGRES_HOST} -U $${POSTGRES_USER} -d $${POSTGRES_DB} -f ./db/data/insert.sql'

clear_data: ## Clear data from the database
	@bash -c 'source .env.local && PGPASSWORD=$${POSTGRES_PASSWORD} psql -h $${POSTGRES_HOST} -U $${POSTGRES_USER} -d $${POSTGRES_DB} -f ./db/data/clear.sql'

run: ## Run the application
	@ECHO_MODE=local go run ./cmd/my_authn_authz/main.go

view_build: ## Build the view files
	@cd ./view && yarn build

help: ## Show Makefile options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'