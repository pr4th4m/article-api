#### For developers
.PHONY: pre
pre:
	set -e; \
    docker pull elasticsearch:7.0.1; \
    docker run -d --name elasticsearch -p 9200:9200 -e "discovery.type=single-node" elasticsearch:7.0.1;

.PHONY: test
test:
	set -e; \
	go test ./...

.PHONY: build
build:
	set -e; \
	go build -o server;
