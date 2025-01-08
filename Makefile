start:
	- go run cmd/main.go
test:
	- go test -v ./...
build:
	- go build -o bin/main cmd/main.go