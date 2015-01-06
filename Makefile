all: build test lint

build:
	go build ./...

test:
	go test ./...

lint:
	golint ./...
	gofmt -w -s . ./example*
	goimports -w . ./example*

clean:
	rm -f *~ ./example*/*~
	rm -f ./example1/example1
	rm -f ./example2/example2
	go clean ./...