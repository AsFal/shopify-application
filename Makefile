.PHONY: build

build:
	- go build -o bin/api ./cmd/api/main.go

unit_tests:
	- go test ./internal/...

clean:
	- rm -r bin

