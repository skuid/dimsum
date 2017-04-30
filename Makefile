REPO=dimsum
ORGANIZATION=095427547185.dkr.ecr.us-west-2.amazonaws.com/skuid

.PHONY: build


build:
	docker run --rm -v $$(pwd):/go/src/github.com/skuid/$(REPO) -w /go/src/github.com/skuid/$(REPO) golang:1.8  go build -v -a -tags netgo -installsuffix netgo -ldflags '-w'
	docker build -t $(ORGANIZATION)/$(REPO) .

clean:
	rm ./$(REPO)


