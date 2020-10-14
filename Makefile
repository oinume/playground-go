all: setup

.PHONY: setup
setup:
	go list -tags=tools -f='{{ join .Imports "\n" }}' ./tools/tools.go | tr -d [ | tr -d ] | xargs -I{} go install {}

.PHONY: lint
lint:
	golangci-lint run -j 4 --out-format=line-number ./...

.PHONY: moq
moq:
	rm -f ./moq/github/*.moq.go
	go generate ./moq/github