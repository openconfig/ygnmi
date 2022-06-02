.PHONY: test integration_test clean

build:
	go build -v ./...

# TODO: Add race tests for specific packages with concurrency.
test:
	go test -v ./...

gen:
	internal/exampleoc/gen.sh

clean:
	find internal/exampleoc -mindepth 1 -maxdepth 1 ! -name gen.go ! -name gen.sh -exec rm -r {} \+
