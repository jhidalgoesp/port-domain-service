lint:
	golangci-lint --config=.golangci_lint.yaml run

test:
	go test ./... -v -count=1

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

run:
	go run ./cmd

docker-build:
	docker build -t ports-app .

docker-run:
	docker run -p 8080:8080 ports-app

