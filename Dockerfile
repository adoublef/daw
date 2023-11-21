# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21
ARG ALPINE_VERSION=3.18

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS baseline
WORKDIR /usr/src

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# required for go-sqlite3
RUN apk add --no-cache --upgrade --latest gcc musl-dev

COPY go.* .
RUN go mod download

COPY . .

FROM baseline AS testing

RUN go test -v -count=1 ./...

FROM baseline AS build

RUN go build \
    -ldflags "-s -w -extldflags '-static'" \
    -buildvcs=false \
    -o /usr/local/bin/ ./...

FROM alpine:${ALPINE_VERSION} AS runtime
WORKDIR /opt

RUN addgroup -g 10001 daw \
    && adduser -G daw -u 10001 daw -D

COPY --from=build /usr/local/bin/daw ./a

USER daw:daw

ENTRYPOINT ["./a"]
CMD ["serve"]