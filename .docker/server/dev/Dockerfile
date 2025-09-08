FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache git bash

# hot-reloading with air
RUN go install github.com/air-verse/air@latest

# golangci-lint for IDE integration (Auto-completion, linting, formatting, etc.)
RUN go install golang.org/x/tools/gopls@latest

COPY go.mod ./
RUN go mod download

EXPOSE 7777

ENTRYPOINT ["air", "-c", ".air.dev.toml"]
