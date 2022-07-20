all: snapshot

dev:
	go run --tags dev . server --http

snapshot: 
	goreleaser release --snapshot --rm-dist
