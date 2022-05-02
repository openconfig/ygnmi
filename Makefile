build:
	go build -v ./...

test:
	go test -v ./...

clean:
	rm -rf test/device
	rm -f test/telemetry.go