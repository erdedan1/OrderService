MIGRATE_URL=postgres://postgres:postgres@postgres:5432/order_db?sslmode=disable

migrate-up:
	migrate -path migrations -database "$(MIGRATE_URL)" up

migrate-down:
	migrate -path migrations -database "$(MIGRATE_URL)" down 1
