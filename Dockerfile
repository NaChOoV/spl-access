# Use the official Golang image as the base image
FROM golang:1.23.5-alpine AS build
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o main main.go

# Production
FROM alpine:latest
WORKDIR /app
COPY --from=build ./app/main .
COPY .env .env

EXPOSE 30002

CMD ["./main"]