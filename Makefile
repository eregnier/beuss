test:
	@export DEBUG=0 && go test

test-trace:
	@clear && export DEBUG=1 && go test

build:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o beuss-server
	
dev:
	@go run main.go

run: build
	./beuss-server

docker-run: build
	docker build . --tag beuss
	docker run -it  -p 6552:6552 beuss