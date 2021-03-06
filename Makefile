PACKAGES := github.com/hekonsek/awsom/aws github.com/hekonsek/awsom github.com/hekonsek/awsom/main
VERSION := 0.3.0

all: format rice silent-test build

rice:
	rice embed-go

build:
	GO111MODULE=on go build -o awsom main/*.go

test: build
	GO111MODULE=on go test -v $(PACKAGES)

silent-test:
	GO111MODULE=on go test $(PACKAGES)

format:
	GO111MODULE=on go fmt $(PACKAGES)

docker-build: build
	docker build -t hekonsek/awsom:$(VERSION) .

docker-push: docker-build
	docker push hekonsek/awsom:$(VERSION)

release: docker-push
	git tag v$(VERSION) && git push --tags

lint: format
	~/go/bin/golint $(PACKAGES)