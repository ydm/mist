.PHONY: build
build:
	go build cmd/mist.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test
