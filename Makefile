VERSION     := $(shell git tag -l --sort=v:refname "v*" | tail -1)-dev
COMMIT      := $(shell git log -n 1 --pretty=format:%h --abbrev=8 2> /dev/null)
DATE        := $(shell date "+%Y-%m-%dT%H:%M:%S+09:00")
BUILD_FLAGS := -ldflags "-s -w -X github.com/sugy/filequeue/cmd.version=$(VERSION) -X github.com/sugy/filequeue/cmd.commit=$(COMMIT) -X github.com/sugy/filequeue/cmd.date=$(DATE)"

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
	go build $(BUILD_FLAGS)
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o filequeue_darwin_arm64
