lint:
	golangci-lint run

test:
	go test  ./...

build:
	go build ./cmd/dump-bakusai
