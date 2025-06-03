.PHONY: clean test

toml2json: go.* *.go
	go build -o $@ ./cmd/toml2json

clean:
	rm -rf toml2json dist/

test:
	go test -v ./...

install:
	go install github.com/fujiwara/toml2json/cmd/toml2json

dist:
	goreleaser build --snapshot --clean
