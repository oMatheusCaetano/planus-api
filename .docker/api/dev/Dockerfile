FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache git bash

# hot-reloading with air
RUN go install github.com/air-verse/air@latest

# database migration tool
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# golangci-lint for IDE integration (Auto-completion, linting, formatting, etc.)
RUN go install golang.org/x/tools/gopls@latest

# Swagger packages for API documentation
RUN wget -q https://github.com/go-swagger/go-swagger/releases/download/v0.30.5/swagger_linux_amd64 \
    -O /usr/local/bin/swagger && chmod +x /usr/local/bin/swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod ./
RUN go mod download

EXPOSE 7777

ENTRYPOINT ["air", "-c", ".air.dev.toml"]
