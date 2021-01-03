.PHONY: generate test-format test-unit test-e2e gosec docker-build docker-publish-image docker-tests dockerize install-tools

generate:
	go generate ./...

test-format:
	go fmt $(go list ./...)
	go vet $(go list ./...)

test-unit: generate
	go test --race -v ./...

test-e2e: generate
	go test -race ./e2e -tags=e2

gosec:
	gosec ./...

docker-build:
	docker build -t pg2000/artifacthub-resource:latest .

docker-publish-image: docker-build
	docker push pg2000/artifacthub-resource:latest

docker-tests:
	docker build --target tests -t artifacthub-resource-tests .

install-tools:
	go list -f '{{range .Imports}}{{.}} {{end}}' ./tools/tools.go | xargs go install