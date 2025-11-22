run:
	go run main.go

build:
	go build -o bin/bcu main.go

test:
	go test ./...

release-snapshot: test
	goreleaser release --snapshot --clean

release: test
	goreleaser release --clean