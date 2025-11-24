run:
	go run main.go

build:
	go build -o bin/bcu main.go

test:
	go test ./...

lint:
	golangci-lint run

check: test lint

release-snapshot: check
	goreleaser release --snapshot --clean

release: check
	goreleaser release --clean