FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . . 

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /usr/local/bin/app ./cmd/server

FROM alpine:latest

WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/app .

ENTRYPOINT ["/usr/local/bin/app"]

EXPOSE 8080