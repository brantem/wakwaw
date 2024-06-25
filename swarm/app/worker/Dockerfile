FROM golang:alpine AS builder
WORKDIR /builder

COPY /shared /shared
WORKDIR /worker
COPY /worker .
RUN go mod download
RUN CGO_ENABLED=0 go build -o app -ldflags="-s -w" example.com/worker

FROM alpine:latest

RUN apk --no-cache add curl
COPY --from=builder /worker/app .
CMD ["./app"]
