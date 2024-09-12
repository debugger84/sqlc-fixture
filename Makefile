.PHONY: build test

build:
	go build ./...

test: bin/sqlc-fixture.wasm
	go test ./...

all: bin/sqlc-fixture bin/sqlc-fixture.wasm

bin/sqlc-fixture: bin go.mod go.sum $(wildcard **/*.go)
	cd plugin && go build -o ../bin/sqlc-fixture ./main.go

bin/sqlc-fixture.wasm: bin/sqlc-fixture
	cd plugin && GOOS=wasip1 GOARCH=wasm go build -o ../bin/sqlc-fixture.wasm main.go

bin:
	mkdir -p bin
