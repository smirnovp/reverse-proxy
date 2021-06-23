.PHONY: run
run:
	go run --race cmd/proxyserver/main.go

.DEFAULT_GOAL := run