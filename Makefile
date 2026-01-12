.PHONY: help createdb dropdb listdb start stop migrate migrateup migratedown

# Default PostgreSQL container name (override with: make createdb CONTAINER=your-container-name)
CONTAINER ?= postgres_container

help: ## Show available commands
	@echo "Available targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

createdb: ## Create the database
	docker exec -it $(CONTAINER) createdb -U postgres m_webapp_go || \
	docker exec -it $(CONTAINER) psql -U postgres -c "CREATE DATABASE m_webapp_go;"

dropdb: ## Drop the database
	docker exec -it $(CONTAINER) dropdb -U postgres m_webapp_go || \
	docker exec -it $(CONTAINER) psql -U postgres -c "DROP DATABASE m_webapp_go;"

listdb: ## List all databases
	docker exec -it $(CONTAINER) psql -U postgres -l

start: ## Start the application
	@bash run.sh &

stop: ## Stop the application
	@pkill -f "modern-web-app" 2>/dev/null && echo "Application stopped" || echo "Application is not running"

migrate: ## Create a new migration file
	migrate create -ext sql -dir db/migrate -seq init_schema

migrateup: ## Run database migrations up
	migrate -path db/migrate -database "postgresql://postgres:Cluster@1@localhost:5432/m_webapp_go?sslmode=disable" -verbose up

migratedown: ## Run database migrations down
	migrate -path db/migrate -database "postgresql://postgres:Cluster@1@localhost:5432/m_webapp_go?sslmode=disable" -verbose down
