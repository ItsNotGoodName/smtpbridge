all: npm snapshot

npm:
	npm i --prefix left/web

dev-backend:
	go run --tags dev . server

dev-frontend:
	npm run dev --prefix left/web

snapshot: 
	goreleaser release --snapshot --rm-dist
