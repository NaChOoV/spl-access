FROM golang:1.23.4-alpine AS build
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "aarch64" ]; then \
    CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o main main.go; \
    else \
    go build -ldflags="-w -s" -o main main.go; \
    fi

# Production
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/main .
COPY --from=build /app/src/repository/sql/*.sql /app/src/repository/sql/

EXPOSE 8000

CMD ["./main"]