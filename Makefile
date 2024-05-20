run:
	go run main.go

build:
	go build -o bin/main main.go

release-snapshot:
	goreleaser release --snapshot --clean
