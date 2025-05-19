FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/dns-server
RUN go build -o dns-server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /usr/local/bin
COPY --from=builder /app/cmd/dns-server/dns-server .
EXPOSE 53/udp
EXPOSE 9090/tcp
ENTRYPOINT ["/usr/local/bin/dns-server"]