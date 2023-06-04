format:
	@go fmt ./...

test:
	@go test -cover -v ./...

cover:
	@go test ./... -coverprofile cover.out
	@go tool cover -html cover.out -o cover.html

run:
	@go run main.go

build:
	@go build -v ./...

install:
	@go install -v ./...

vulncheck:
	@govulncheck ./...

pre-commit:
	@pip3 install pre-commit
	@pre-commit install
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.1
