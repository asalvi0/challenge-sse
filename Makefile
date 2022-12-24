default:
	@CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o . ./cmd/server.go

.PHONY: test
test:
	@go test ./...
