.PHONY: test integration_test clean

build:
	go build -v ./...

test:
	go test -v ./...
	go test -race -v ./...

integration_test:
	go generate ./...
	go build ./test

clean:
	rm -rf test/device
	rm -f test/telemetry.go