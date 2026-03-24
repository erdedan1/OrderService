MIGRATE_URL=postgres://postgres:postgres@postgres:5432/order_db?sslmode=disable

migrate-up:
	goose -dir internal/migrations postgres "$(MIGRATE_URL)" up

migrate-down:
	goose -dir internal/migrations postgres "$(MIGRATE_URL)" down
