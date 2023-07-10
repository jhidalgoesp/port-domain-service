# Port Domain Service

This is an example of a microservice built using hexagonal architecture.

```bash
.
├── Makefile                    <-- Make to automate tasks
├── README.md                   <-- This instructions file
├── cmd                         <-- Application entrypoint
├── internal                    
│   ├── database                <-- Memory Database Adapter
│   └── file_reader             <-- File Reader Adapter
│   └── domain                  <-- Application Core/Domain
```
## Design decisions

- The file reader implementation uses `encoding/json` and `json.Decoder` from the standard library, so the JSON file is read in chunks.
- The service implements a handler to execute the ports load and database upsert against any HTTP request/endpoint.
- The implementation File Reader only works with the JSON format that was provided.

## Requirements

* [Golang](https://go.dev) 1.20
* [Golanglint-ci](https://golangci-lint.run/)
* [Docker](https://www.docker.com/)

## Commands

```sh
make run # or go run ./cmd
```

```sh
make test # Run Unit Tests
```

```sh
make cover # Coverage report
```

```sh
make lint # Run golanci-lint
```

Running with Docker

```sh
make docker-build # or docker build -t ports-app .
```

```sh
make docker-run # or docker run -p 8080:8080 ports-app
```