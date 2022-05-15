SHELL := bash
.ONESHELL:
MAKEFLAGS += --no-builtin-rules

NOCACHE := $(if $(NOCACHE),"--no-cache")
MIGRATE_DSN := "postgres://stat_test_task:stat_test_task@postgres:5432/stat_test_task?sslmode=disable"

export APP_NAME := pg-stat-test-task
export DOCKER_REPOSITORY := radon1
export VERSION := $(if $(VERSION),$(VERSION),$(if $(COMMIT_SHA),$(COMMIT_SHA),$(shell git rev-parse --verify HEAD)))

.PHONY: help
help: ## List all available targets with help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: generate
generate: ## run all golang codegen
	@go generate ./...

.PHONY: build
build: build-helper build-prod ## Build all containers

.PHONY: build-helper
build-helper:
	@docker build ${NOCACHE} --pull -f ./build/helper.Dockerfile -t ${DOCKER_REPOSITORY}/${APP_NAME}-helper:${VERSION} .

.PHONY: build-prod
build-prod:
	@docker build ${NOCACHE} --pull -f ./build/Dockerfile -t ${DOCKER_REPOSITORY}/${APP_NAME}:${VERSION} .

.PHONY: migration-up
migration-up: ## Run develop migrations
	@docker-compose run --rm -T helper migrate -verbose -path ./migrations -database ${MIGRATE_DSN} up

.PHONY: migration-down
migration-down: ## Rollback develop migrations
	@docker-compose run --rm -T helper migrate -verbose -path ./migrations -database ${MIGRATE_DSN} down

.PHONY: run-dev-env
run-dev-env:
	@docker-compose up -d postgres

.PHONY: run
run: run-dev-env migration-up ## Run develop docker-compose
	@docker-compose up app

.PHONY: stop
stop: ## Stop all develop containers
	@docker-compose down -v

.PHONY: test-short
test-short: ## Run unit tests
	@go test ./... -cover -short

.PHONY: test-long
test-long: ## Run all tests (unit/integrations)
	@make run-dev-env && make migration-up && make test-long-up; make stop

.PHONY: tests-long-up
test-long-up:
	@docker-compose run --rm helper go test -cover ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

.PHONY: psql-exec
psql-exec: ## exec to psql
	@docker-compose run --rm postgres psql -h postgres -p 5432 -U stat_test_task