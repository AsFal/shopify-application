.PHONY: build

build:
	- go build -o bin/api ./cmd/api/main.go

test:
	- go test ./internal/...

clean:
	- rm -r bin

