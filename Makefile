all: dep-web snapshot

snapshot:
	goreleaser release --snapshot --clean

dev:
	air

dev-web:
	cd web && pnpm run dev

clean:
	rm -rf smtpbridge_data

dep: dep-air dep-sqlc

dep-air:
	go install github.com/cosmtrek/air@latest

dep-sqlc:
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

dep-web:
	cd web && pnpm install
