FROM golang:1.27.4-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
# COPY vendor .
COPY .. .
RUN mkdir bin
RUN go build -o bin ./cmd/server

FROM alpine:3.24
RUN apk add --no-cache ca-certificates
WORKDIR /bin
COPY --from=builder /build/bin/server server
ENTRYPOINT ["/bin/server"]