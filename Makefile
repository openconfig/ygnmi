.PHONY: test integration_test clean

build:
	go build -v ./...

# TODO: Add race tests for specific packages with concurrency.
test:
	go test -v ./...

integration_test:
	test/gen.sh
	go build ./test

clean:
	rm -rf test/device
	rm -f test/telemetry.go
