clean:
	go clean ./

build:
	go build -o ./irishub-load ./

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./irishub-load ./

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs goimports -w -local github.com/irisnet/irishub-load
