.PHONY: build

build:
	- go build -o bin/api ./cmd/api/main.go

unit_tests:
	- go test ./internal/...

system_tests:
	- go test ./tests/... -tags system

clean:
	- rm -r bin

