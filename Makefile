
build:  ## build and lint
	go build ./...
	gometalinter \
                --vendor \
                --vendored-linters \
                --deadline=60s \
                --disable-all \
                --enable=goimports \
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
test:  ## just test
	go test -cover .

clean: ## cleanup
	rm -f ./example1/example1
	rm -f ./example2/example2
	go clean ./...
	git gc

# https://www.client9.com/self-documenting-makefiles/
help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
        printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
        }' $(MAKEFILE_LIST)
.DEFAULT_GOAL=help
.PHONY=help
