run:
	go run main.go

build:
	go build -o bin/bcu main.go

release-snapshot:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean