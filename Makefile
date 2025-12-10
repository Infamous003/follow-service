.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/n] ' && read ans && [ $${ans:-n} = y ]

## run/api: run the cmd/api application
.PHONY: run
run:
	@go run ./cmd/ 

.PHONY: db/psql
db/psql:
	@psql "$(FOLLOW_DB_DSN)"

.PHONY: db/migrations/new
db/migrations/new:
	@echo "Creating new migration: $(name)"
	@migrate create -ext sql -dir ./migrations -seq $(name)

.PHONY: db/migrations/up
db/migrations/up:
	@echo "Applying up migrations to database"
	@migrate -path ./migrations -database "$(FOLLOW_DB_DSN)" up

.PHONY: db/migrations/down
db/migrations/down:
	@echo "Applying down migrations to database by 1"
	@migrate -path ./migrations -database "$(FOLLOW_DB_DSN)" down 1