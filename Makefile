export DB_DIR=./smtpbridge_data
export DB_FILE=smtpbridge.db
export DB_PATH="$(DB_DIR)/$(DB_FILE)"
export DEV_IP=127.0.0.1

-include .env

snapshot:
	goreleaser release --snapshot --clean

run:
	go run ./cmd/smtpbridge

preview:
	cd web && pnpm run build && cd .. && go run ./cmd/smtpbridge

clean:
	rm -rf "$(DB_DIR)" && mkdir "$(DB_DIR)"

gen: db-migrate gen-jet gen-templ

tooling: tooling-air tooling-jet tooling-goose tooling-atlas tooling-templ tooling-goreleaser

# Development

dev:
	air

dev-web:
	cd web && pnpm install && pnpm run dev

# Database

db-ui:
	podman run -it --rm \
		-p 8090:8080 \
		-v "$(DB_DIR):/data" \
		-e "SQLITE_DATABASE=$(DB_FILE)" \
		docker.io/coleifer/sqlite-web

db-inspect:
	atlas schema inspect --env local

db-migration:
	atlas migrate diff $(name) --env local

db-migrate:
	goose -dir migrations/sql sqlite3 "$(DB_PATH)" up

# Generation

gen-jet:
	jet -source=sqlite -dsn="$(DB_PATH)" -path=./internal/jet -ignore-tables goose_db_version,_dummy
	rm -rf ./internal/jet/model

gen-templ:
	cd web && templ generate

# Tooling

tooling-air:
	go install github.com/air-verse/air@latest

tooling-jet:
	go install github.com/go-jet/jet/v2/cmd/jet@latest

tooling-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

tooling-atlas:
	go install ariga.io/atlas/cmd/atlas@latest

tooling-templ:
	go install github.com/a-h/templ/cmd/templ@latest

tooling-goreleaser:
	go install github.com/goreleaser/goreleaser@latest
