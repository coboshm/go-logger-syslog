SOURCES          ?= $(shell find . -name "*.go" | grep -v vendor/)
PACKAGES         ?= $(shell go list ./... | grep -v vendor)
GOTOOLS          ?= github.com/GeertJohan/fgt \
					golang.org/x/tools/cmd/goimports \
					github.com/golang/lint/golint \
					github.com/golang/mock/gomock \
					github.com/golang/mock/mockgen \
					github.com/kisielk/errcheck \
					honnef.co/go/tools/cmd/staticcheck

deps: tools
	go get -t ./...

test: deps
	go test $(PACKAGES)

test-verbose: deps
	go test -v $(PACKAGES)

lint: tools
	fgt go fmt $(PACKAGES)
	fgt goimports -w $(SOURCES)
	echo $(PACKAGES) | xargs -L1 fgt golint
	fgt go vet $(PACKAGES)
	fgt errcheck -ignore Close $(PACKAGES)
	staticcheck $(PACKAGES)
.SILENT: lint

check: test lint

tools:
	go get $(GOTOOLS)
.SILENT: tools

tools-update:
	go get -u $(GOTOOLS)

build: deps
	mkdir -p build
	go build -o app-logger-test
.PHONY: build