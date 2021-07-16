.PHONY: build clean fmt vet test bench cover profiling

ARGS=
COVERPROFILE=coverage.txt
DEBUG=

build: clean fmt vet test
	go build internal/web/...
	go build cmd/web/...
	go build cmd/worker/...

clean:
	rm -rf bin/
	if [ -f "cmd/main" ] ; then rm cmd/main; fi
	if [ -f "coverage.txt" ] ; then rm coverage.txt; fi

fmt:
	go fmt ./cmd/...
	go fmt ./internal/...
	go fmt ./pkg/...

vet:
	go fmt ./cmd/...
	go fmt ./internal/...
	go fmt ./pkg/...

test:
	go test ./internal/... -coverprofile=$(COVERPROFILE) -covermode atomic -v -race $(DEBUG) $(ARGS)

cover:
	$(eval COVERPREFILE += -coverprofile=coverage.out)
	go test ./internal/... -cover $(COVERPREFILE) -race $(ARGS) $(DEBUG)
	go tool cover -html=coverage.out
	rm -f coverage.out

profiling:
	go test -bench=/internal -cpuprofile cpu.out -memprofile mem.out $ARGS