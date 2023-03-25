build-server:
	@go build -o bin/app .

server: build-server
	@bin/app

test-server:
	@go test -v ./...