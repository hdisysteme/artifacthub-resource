# stage: builder
FROM golang:1.15-buster AS builder

WORKDIR /concourse/concourse-resource
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

ENV CGO_ENABLED 0

RUN go build -o /assets/check github.com/PG2000/artifacthub-resource/cmd/check
RUN go build -o /assets/in github.com/PG2000/artifacthub-resource/cmd/in

# stage: tests
FROM builder as tests
WORKDIR /app
COPY --from=builder /concourse/concourse-resource /app
ENV CGO_ENABLED 1
RUN set -e; \
    make generate && \
    make test-format && \
    make test-unit && \
    make test-e2e && \
    make gosec

# stage: resource
FROM ubuntu:focal AS resource

RUN apt-get update && apt-get upgrade -y --no-install-recommends && apt-get install -y --no-install-recommends \
    ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder /assets /opt/resource

# final output stage
FROM resource