MAIN_PACKAGE_PATH := ./cmd/raffle
BINARY_NAME := raffle

.PHONY: default
default: build

.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

.PHONY: test
test:
	go test -v -race -buildvcs ./...

.PHONY: build
build: tidy test
	go build -o=$(CURDIR)/dist/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

.PHONY: clean
clean:
	@rm -rf $(CURDIR)/dist
