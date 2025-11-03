ARG GO_VERSION=1.25
ARG DEBIAN_VERSION=bookworm-20251020-slim

FROM golang:${GO_VERSION} AS build

RUN apt-get update && apt-get install -y git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -ldflags="-w -s -X main.version=${VERSION:-dev}" -a -o /app/backend .

FROM debian:${DEBIAN_VERSION} AS runtime

RUN useradd -u 1000 -s /bin/nologin app 
USER app

WORKDIR /app
RUN mkdir -p /app/data

COPY --from=build --chown=app:app --chmod=755 /app/backend /app/

EXPOSE 8080

ENTRYPOINT ["/app/backend"]