all: setup

.PHONY: setup
setup:
	cd tools && go list -f='{{ join .Imports "\n" }}' . | tr -d [ | tr -d ] | xargs -I{} go install {}

.PHONY: lint
lint:
	golangci-lint run -j 4 --out-format=line-number ./...

.PHONY: moq
moq:
	rm -f ./moq/github/*.moq.go
	go generate ./moq/github