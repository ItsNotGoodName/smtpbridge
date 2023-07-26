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

dep: dep-air dep-sqlc dep-web

dep-air:
	go install github.com/cosmtrek/air@latest

dep-sqlc:
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

dep-web:
	cd web && pnpm install
