.PHONY: build
build:
	go build -v ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: test-static-files
test-static-files:
	GENERATE_STATIC_FILES=true go test -v ./...
