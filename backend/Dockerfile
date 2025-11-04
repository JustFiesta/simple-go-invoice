ARG GO_VERSION=1.25
ARG DEBIAN_VERSION=bookworm-20251020-slim

ARG APP_VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

FROM golang:${GO_VERSION} AS build

RUN apt-get update && apt-get install -y git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -ldflags="-w -s \
    -X main.Version=${APP_VERSION} \
    -X main.Commit=${GIT_COMMIT} \
    -X main.BuildDate=${BUILD_DATE}" \
    -a -o /app/backend .

FROM debian:${DEBIAN_VERSION} AS runtime

ARG APP_VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

LABEL version=${APP_VERSION} \
      commit=${GIT_COMMIT} \
      build_date=${BUILD_DATE}

RUN useradd -u 1000 -s /bin/nologin app 
USER app

WORKDIR /app
RUN mkdir -p /app/data

COPY --from=build --chown=app:app --chmod=755 /app/backend /app/

EXPOSE 8080

ENTRYPOINT ["/app/backend"]