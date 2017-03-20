
all:
	go build ./...
	gometalinter \
                --vendor \
                --vendored-linters \
                --deadline=60s \
                --disable-all \
                --enable=goimports \
                --enable=aligncheck \
                --enable=vetshadow \
                --enable=varcheck \
                --enable=structcheck \
                --enable=deadcode \
                --enable=ineffassign \
                --enable=unconvert \
                --enable=goconst \
                --enable=golint \
                --enable=gosimple \
                --enable=gofmt \
                --enable=misspell \
                --enable=staticcheck \
                .
	go test -cover .

clean:
	rm -f ./example1/example1
	rm -f ./example2/example2
	go clean ./...
	git gc

ci: build test lint

docker-ci:
	docker run --rm \
		-v $(PWD):/go/src/github.com/client9/reopen \
		-w /go/src/github.com/client9/reopen \
		nickg/golang-dev-docker \
		make ci

.PHONY: ci docker-ci
