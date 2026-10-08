.PHONY: test test-integration test-e2e lint sqlc run

test: ## Unit test (domain, app, transport)
	go test -race ./...

test-integration: ## Repository với Postgres thật (cần Docker)
	go test -tags=integration ./internal/...

test-e2e: ## E2E (cần Docker)
	go test -tags=e2e ./test/e2e/...

lint: ## golangci-lint (có depguard)
	golangci-lint run ./...

sqlc: ## Sinh lại code từ queries.sql (không sửa tay thư mục db/)
	sqlc generate

run: ## Chạy local với adapter in-memory
	go run ./cmd/server
