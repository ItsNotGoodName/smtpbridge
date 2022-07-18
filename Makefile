all: snapshot

snapshot: 
	goreleaser release --snapshot --rm-dist
