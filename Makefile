start:
	- go run cmd/main.go
test:
	- go test -v ./...
build:
	- go build -o bin/main cmd/main.go
swagger:
	- swag init  --dir cmd,api --parseDependency true