.DEFAULT_GOAL := run

http_path := ./cmd/http
pkgs_path := ./pkgs/*
bins := ./http

fmt:
	go fmt $(http_path)
	go fmt $(pkgs_path)
.PHONY:fmt

lint: fmt
	golint $(http_path)
	golint $(pkgs_path)
.PHONY:lint

vet: lint
	go vet $(http_path)
	go vet $(pkgs_path)
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
