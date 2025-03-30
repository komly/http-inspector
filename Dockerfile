FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY http-inspector.go .
RUN go build -o http-inspector http-inspector.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/http-inspector .
EXPOSE 8080
CMD ["./http-inspector"] 