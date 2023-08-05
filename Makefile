gen:
	jet -source=sqlite -dsn="./smtpbridge_data/smtpbridge.db" -path=./internal/dbgen -ignore-tables bun_migrations,bun_migration_locks

snapshot:
	goreleaser release --snapshot --clean

preview:
	cd web && pnpm run build && cd .. && go run .

dev:
	air

dev-web:
	cd web && pnpm run dev

clean:
	rm -rf smtpbridge_data

dep: dep-air dep-jet dep-web

dep-air:
	go install github.com/cosmtrek/air@latest

dep-jet:
	go install github.com/go-jet/jet/v2/cmd/jet@latest

dep-web:
	cd web && pnpm install
