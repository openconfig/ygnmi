.PHONY: test integration_test clean

build:
	go build -v ./...

# TODO: Add race tests for specific packages with concurrency.
test:
	go test -v ./...

gen:
	internal/exampleoc/gen.sh

clean:
	rm -rf test/device
	rm -f test/telemetry.go
