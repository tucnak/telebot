linter:
	./bin/golangci-lint run ./...
.PHONY: linter

linter.download:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
.PHONY: linter.download

test:
	env CGO_ENABLED=1 go test -race ./...
.PHONY: test

alltest: linter test
.PHONY: alltest
