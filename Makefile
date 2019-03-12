PACKAGES := github.com/hekonsek/awsom github.com/hekonsek/awsom/cmd github.com/hekonsek/awsom/main

all: format rice silent-test build

rice:
	rice embed-go

build:
	GO111MODULE=on go build main/awsom.go

test: build
	GO111MODULE=on go test -v $(PACKAGES)

silent-test:
	GO111MODULE=on go test $(PACKAGES)

format:
	GO111MODULE=on go fmt $(PACKAGES)

docker-build: build
	docker build -t hekonsek/awsom .

docker-push: docker-build
	docker push hekonsek/awsom

lint: format
	~/go/bin/golint $(PACKAGES)