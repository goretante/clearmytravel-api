DATABASE_URL ?= postgres://clearmytravel:clearmytravel@localhost:5432/clearmytravel?sslmode=disable

.PHONY: db-up db-down db-migrate db-seed db-reset test

db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down

db-migrate:
	DATABASE_URL=$(DATABASE_URL) goose -dir migrations postgres "$(DATABASE_URL)" up

db-seed:
	docker exec -i clearmytravel-postgres \
		psql -U clearmytravel -d clearmytravel < seed/development.sql

db-reset:
	docker-compose down -v
	docker-compose up -d postgres
	sleep 2
	DATABASE_URL=$(DATABASE_URL) goose -dir migrations postgres "$(DATABASE_URL)" up
	docker exec -i clearmytravel-postgres \
		psql -U clearmytravel -d clearmytravel < seed/development.sql

test:
	go test ./...