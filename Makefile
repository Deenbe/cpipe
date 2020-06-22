.PHONY: build clean install test

TARGET=
SUFFIX=$(GOOS)_$(GOARCH)
ifneq ("${SUFFIX}", "_")
TARGET=_$(SUFFIX)
endif

AWS_REGION=
ifeq ("${AWS_REGION}", "")
AWS_REGION="us-west-1"
endif

build:
	CGO_ENABLED=0 go build -o "./build/cpipe${TARGET}${EXTENSION}"

clean:
	rm -rf ./build

test: build
	CGO_ENABLED=0 AWS_REGION="${AWS_REGION}" go test -timeout 5s -covermode=count -coverpkg="cpipe" -coverprofile=build/cover.out -tags=integration ./...

install: build
	cp "./build/cpipe${TARGET}" /usr/local/bin/
