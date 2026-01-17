-include .env
export 

print-vars:
	@echo $(DB_URL)
	@echo $(MIGRATIONS_DIR)
	@echo $(DRIVER)

postgres:
	docker run --name simple-bank -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=derick -e  -p 5432:5432 -d postgres:18.1-alpine3.23

create-db:
	docker exec -it simple-bank createdb --username=admin --owner=admin simple_bank

drop-db:
	docker exec -it simple-bank dropdb --username=admin --owner=admin simple_bank

goose-create-file:
	@read -p "File name: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql

migrate-up:
	goose -dir $(MIGRATIONS_DIR) $(DRIVER) $(DB_URL) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) $(DRIVER) $(DB_URL) down

sqlc:
	sqlc generate

test:
	go test -v -cover ./internals/models/

build-tags:
	./scripts/build-tags.sh "$(PR_TITLE)"

.PHONY: postgres create-db drop-db goose-create-file migrate-up migrate-down sqlc test
