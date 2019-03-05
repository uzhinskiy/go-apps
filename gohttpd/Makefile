GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=./src/$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")
PID=/tmp/go-$(GONAME).pid
GOOS=Linux

all :: get build


get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get -d $(GOFILES)

build:
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -ldflags "-s -w" -o bin/$(GONAME) $(GOFILES)

install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

run:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES)

clear:
	@clear

clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY:	build get install run clean

