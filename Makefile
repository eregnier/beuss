test:
	@export DEBUG=0 && go test

test-trace:
	@clear && export DEBUG=1 && go test

build:
	go build main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main
	
docker-run: build
	docker build . --tag beuss
	docker run -it  -p 6552:6552 beuss