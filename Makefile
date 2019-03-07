PACKAGES := github.com/hekonsek/awsom

all: format rice build silent-test

rice:
	rice embed-go

build:
	GO111MODULE=on go build $(PACKAGES)

test: build
	GO111MODULE=on go test -v $(PACKAGES)

silent-test:
	GO111MODULE=on go test $(PACKAGES)

format:
	GO111MODULE=on go fmt $(PACKAGES)