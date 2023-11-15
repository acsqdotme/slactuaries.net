.DEFAULT_GOAL := run

http_path := ./cmd/http
bins := ./http

fmt:
	go fmt $(http_path)
.PHONY:fmt

lint: fmt
	golint $(http_path)
.PHONY:lint

vet: lint
	go vet $(http_path)
.PHONY:vet

run: vet
	go run $(http_path)
.PHONY:run

build:
	go build $(http_path)
.PHONY:build

clean:
	rm -fv $(bins)
.PHONY:clean
