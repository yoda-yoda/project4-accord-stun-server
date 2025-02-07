# Build Stage
FROM golang:1.23.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o stunserver .

# Runtime Stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/stunserver .

EXPOSE 3479/udp

ENTRYPOINT ["./stunserver"]
