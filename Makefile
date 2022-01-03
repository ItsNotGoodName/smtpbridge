all: npm snapshot

npm:
	npm i --prefix left/web

dev-backend:
	air

dev-frontend:
	npm run css-watch --prefix left/web

snapshot: 
	goreleaser release --snapshot --rm-dist
