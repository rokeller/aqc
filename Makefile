VERSION_RAW = $(shell git describe --tags 2>/dev/null || echo '0.0.0')

aqc: $(wildcard **/*.go)
	@go build -ldflags \
		"-X github.com/rokeller/aqc/cmd.version=${VERSION_RAW}-local"

.PHONY: release
release:
	@go build -ldflags "-s -w \
		-X github.com/rokeller/aqc/cmd.version=${VERSION_RAW}-local"

.PHONY: test
test: aqc
	@go test ./...

.PHONY: cover
cover: aqc
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go-cover-treemap -coverprofile coverage.out > coverage.svg

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: update
update:
	@go get -u ./...

.PHONY: clean
clean:
	@rm -rf aqc
