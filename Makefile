.PHONY: build clean fmt vet test bench cover profiling

ARGS=-args -httpPort=":3333"
COVERPROFILE=coverage.txt
DEBUG=DEBUG

build: clean fmt vet test
	go build internal/web/...
	go build cmd/web/...
	go build cmd/worker/...

clean:
	rm -rf bin/
	if [ -f "cmd/main" ] ; then rm cmd/main; fi
	if [ -f "coverage.txt" ] ; then rm coverage.txt; fi

fmt:
	go fmt ./cmd/web/*
	go fmt ./cmd/practice/*
	go fmt ./cmd/worker/*
	go fmt ./internal/web/*
	go fmt ./pkg/zerorpc/*

vet:
	go fmt ./cmd/web/*
	go fmt ./cmd/worker/*
	go fmt ./internal/web/*
	go fmt ./pkg/zerorpc/*

test:
	go test ./internal/**/* -coverprofile=$(COVERPROFILE) -race $(DEBUG) $(ARGS) -covermode=automic -v
