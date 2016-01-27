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
	git gc

ci: build test lint

docker-ci:
	docker run --rm \
		-e COVERALLS_REPO_TOKEN=$(COVERALLS_REPO_TOKEN) \
		-v $(PWD):/go/src/github.com/client9/reopen \
		-w /go/src/github.com/client9/reopen \
		nickg/golang-dev-docker \
		make ci

.PHONY: ci docker-ci
