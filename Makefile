NPM_PREFIX := podman run --rm -it -v "$(shell pwd)/left/http/static:/work" -w /work docker.io/library/node:16

all: npm snapshot

npm:
	$(NPM_PREFIX) npm install

npm-login:
	$(NPM_PREFIX) bash

npm-dev:
	$(NPM_PREFIX) npm run dev

npm-build:
	$(NPM_PREFIX) npm run build

dev:
	go run --tags dev . --watch --memory

snapshot: npm-build
	goreleaser release --snapshot --rm-dist
