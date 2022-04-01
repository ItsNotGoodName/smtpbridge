NPM_PREFIX := podman run --rm -it -p 3000:3000 -v "$(shell pwd)/left/web:/work" -w /work docker.io/library/node:16

all: npm build-frontend snapshot

npm:
	$(NPM_PREFIX) npm install

npm-login:
	$(NPM_PREFIX) bash

dev-backend:
	go run --tags dev . server

dev-frontend:
	$(NPM_PREFIX) npm run dev

build-frontend:
	$(NPM_PREFIX) npm run build

snapshot: 
	goreleaser release --snapshot --rm-dist
