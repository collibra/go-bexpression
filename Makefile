# To try different version of Go
GO := go

gotestsum := go run gotest.tools/gotestsum@latest

generate-mocks:
	$(GO) run github.com/vektra/mockery/v3

generate: generate-mocks
	go generate ./...

build: generate
	go build ./...

test:
	$(gotestsum) --debug --format testname -- -race -mod=readonly -v ./...

test-coverage:
	$(gotestsum) --debug --format testname -- -race -mod=readonly -v -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html

lint:
	golangci-lint run ./...
	go fmt ./...