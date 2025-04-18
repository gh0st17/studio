FROM golang:1.24.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w" -o studio .
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/studio .
COPY web/html /app/web/html
ENTRYPOINT ["./studio"]