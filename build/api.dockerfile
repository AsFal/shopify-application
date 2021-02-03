FROM golang:alpine AS builder

RUN apk add -U --no-cache ca-certificates
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o api cmd/api/main.go

# Development Image
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/api /
EXPOSE $PORT
ENTRYPOINT ["/api"]