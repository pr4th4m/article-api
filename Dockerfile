# Builder image
FROM golang:1.12.3-alpine3.9 as builder

# Update and install git
RUN set -x \
    && apk update --quiet \
    && apk add --no-cache \
      git

WORKDIR /usr/src/server
COPY . .
RUN GOOS=linux go build -ldflags="-w -s" -v -o /go/bin/server github.com/pratz/nine-article-api


# Deployer image
FROM alpine:3.9

# Copy binary from builder image
COPY --from=builder /go/bin/server /go/bin/server

WORKDIR /go/bin

ENTRYPOINT ["/go/bin/server"]
CMD ["-host", "0.0.0.0", "-loglevel", "debug"]
