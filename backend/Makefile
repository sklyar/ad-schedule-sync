V := @

.PHONY: dependencies docker-up application docker-down start-app unit-tests integration-tests

dependencies:
	$(call print-target)
	$(V)go mod tidy

docker-up:
	$(call print-target)
	$(V)docker-compose up -d
	@echo "Waiting for database to be ready..."
	$(V)while [ $$(docker inspect -f '{{.State.Running}}' postgres-transactor-example) != "true" ]; do sleep 1; done

application: docker-up
	$(call print-target)
	$(V)DATABASE_URL=postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable go run ./main.go

docker-down:
	$(call print-target)
	$(V)docker-compose down

start-app: dependencies application docker-down

unit-tests: dependencies
	@echo "Running unit tests..."
	@go test -v -short ./...

migrate-up:
	$(call print-target)
	$(V)migrate -source ./scripts/migrations -database postgres://localhost:5432/database up

migrate-down:
	$(call print-target)
	$(V)go run ./migrations/migrate.go down 1

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
