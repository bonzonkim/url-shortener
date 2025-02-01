FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum .
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 go build -o url-shortener

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortener .

EXPOSE 8080

CMD ["./url-shortener"]
