FROM golang:alpine AS builder
WORKDIR /builder

COPY /shared /shared
WORKDIR /api
COPY /api .
RUN go mod download
RUN CGO_ENABLED=0 go build -o server -ldflags="-s -w" example.com/api

FROM alpine:latest

RUN apk --no-cache add curl
COPY --from=builder /api/server .
CMD ["./server"]
