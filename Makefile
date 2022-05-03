.PHONY: test integration_test clean

build:
	go build -v ./...

test:
	go test -v ./...
	go test -race -v ./...

integration_test:
	test/gen.sh
	go build ./test

clean:
	rm -rf test/device
	rm -f test/telemetry.go
