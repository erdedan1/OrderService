MIGRATE_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

migrate-up:
	goose -dir internal/migrations postgres "$(MIGRATE_URL)" up

migrate-down:
	goose -dir internal/migrations postgres "$(MIGRATE_URL)" down

migrate-reset:
	goose -dir internal/migrations postgres "$(MIGRATE_URL)" reset

probe:
	go run ./cmd/test2