.PHONY: generate
generate:
	go generate ./...

.PHONY: test-format
test-format:
	go fmt $(go list ./...)
	go vet $(go list ./...)

.PHONY: test-unit
test-unit: generate
	go test --race -v ./...

.PHONY: test-e2e
test-e2e: generate
	go test -race ./e2e -tags=e2

.PHONY: gosec
gosec:
	gosec ./...

.PHONY: docker-build
docker-build:
	docker build -t pg2000/artifacthub-resource:latest .

.PHONY: docker-publish-image
docker-publish-image: docker-build
	docker push pg2000/artifacthub-resource:latest

.PHONY: docker-tests
docker-tests:
	docker build --target tests -t artifacthub-resource-tests .

.PHONY: install-tools
install-tools:
	go list -f '{{range .Imports}}{{.}} {{end}}' ./tools/tools.go | xargs go install
