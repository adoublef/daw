# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21
ARG ALPINE_VERSION=3.18

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS baseline
WORKDIR /usr/src

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# required for go-sqlite3 & nodejs
RUN apk add --no-cache --upgrade --latest gcc musl-dev
RUN apk add nodejs npm

COPY go.* .
RUN go mod download

COPY . .

FROM baseline AS testing

RUN go test -v -count=1 ./...

FROM baseline AS build

RUN go generate ./... \
    && go build \
    -ldflags "-s -w" \
    -buildvcs=false \
    -o /usr/local/bin/ ./...

FROM alpine:${ALPINE_VERSION} AS runtime
WORKDIR /opt

ENV GID=10001
ENV UID=10001
ARG SQL_DIR=/data/sqlite

RUN addgroup -g ${GID} daw \
    && adduser -G daw -u ${UID} daw -D \
    && mkdir -p ${SQL_DIR} \
    && chown -R daw:daw ${SQL_DIR}

COPY --from=build /usr/local/bin/daw ./a

USER daw:daw

ENTRYPOINT ["./a"]
CMD ["serve"]