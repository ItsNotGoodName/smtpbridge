all: npm snapshot

npm:
	cd web && pnpm install

snapshot:
	goreleaser release --snapshot --clean

dev:
	air

dev-vite:
	cd web && pnpm run dev

dev-uno:
	cd web && pnpm run dev-uno

clean:
	rm -rf smtpbridge_data

dep: dep-air dep-sqlc

dep-air:
	go install github.com/cosmtrek/air@latest

dep-sqlc:
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
