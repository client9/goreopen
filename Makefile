all: build test lint

build:
	go build ./...

test:
	go test ./...

lint:
	gofmt -w -s . ./example*
	goimports -w . ./example*

clean:
	rm -f *~ ./example*/*~
	go clean ./...