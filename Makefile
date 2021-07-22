.PHONY: build clean fmt vet test bench cover profiling

ARGS=
COVERPROFILE=coverage.txt
DEBUG=

build: clean fmt vet test
	go build pkg/...
	go build internal/...
	go build cmd/web/...
	go build cmd/worker/...

clean:
	rm -rf bin/
	if [ -f "cmd/main" ] ; then rm cmd/main; fi
	if [ -f "coverage.txt" ] ; then rm coverage.txt; fi

fmt:
	go fmt ./pkg/...
	go fmt ./internal/...
	go fmt ./cmd/...

vet:
	go fmt ./pkg/...
	go fmt ./internal/...
	go fmt ./cmd/...

test:
	go test ./cmd/... -coverprofile=$(COVERPROFILE) -covermode atomic -v -race $(DEBUG) $(ARGS)

cover:
	$(eval COVERPREFILE += -coverprofile=coverage.out)
	go test ./cmd/... -cover $(COVERPREFILE) -race $(ARGS) $(DEBUG)
	go tool cover -html=coverage.out
	rm -f coverage.out

benchLogin:
	go test -v ./cmd/web/router -test.bench Login -test.cpuprofile benchmark_login_cpu.out -test.memprofile benchmark_login_mem.out $ARGS

runRpc:
	go run ./cmd/worker

runWeb:
	go run ./cmd/web


