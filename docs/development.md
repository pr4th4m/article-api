**Development**

- Pre-requisites
    - Docker 17.05 or higher
    - Go 1.12.2 or higher
    - export GO111MODULE=on

- Install dependencies, elasticsearch 7.x is required

        make pre

- Run tests

        make test

- Build code

        make build

- Run server

        ./server -loglevel debug

- Server options

        Usage of ./server:

        -host string
              Server address (default "localhost")
        -loglevel value
              Control log level
        -port string
              Server address (default "8080")

- [Start using API](usage.md)
