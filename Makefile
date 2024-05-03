.PHONY: setup
setup:
	go install golang.org/x/lint/golint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: test
test: setup
	go test -v ./...

.PHONY: lint
lint: setup
	go vet ./...
	golint -set_exit_status ./...
	staticcheck ./...

.PHONY: build
build:
	go build
