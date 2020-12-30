.PHONY: generate test-format test-unit test-e2e dockerize

generate:
	go generate ./...

test-format: generate
	go fmt $(go list ./...)
	go vet $(go list ./...)

test-unit: generate
	go test --race -v ./...

test-e2e: generate
	go test -race ./e2e -tags=e2e

dockerize:
	docker build -t pg2000/artifacthub-resource:latest .
	docker push pg2000/artifacthub-resource:latest
