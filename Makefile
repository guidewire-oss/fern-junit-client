.PHONY: build
build:
	go build -v ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: generate-test-files
generate-test-files:
	go test ./pkg/client/
