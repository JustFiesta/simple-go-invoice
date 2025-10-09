ARG GO_VERSION=1.25

FROM golang:${GO_VERSION}
WORKDIR /app
COPY . .
RUN go mod download
CMD ["go", "run", "main.go"]
