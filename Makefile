clean:
	go clean ./

build:
	go build -o ./irishub-load ./

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./irishub-load ./

