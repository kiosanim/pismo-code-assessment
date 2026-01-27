FROM golang:1.25.6-alpine3.23 as builder
WORKDIR /app
COPY . .
RUN go build -o api ./cmd/api

FROM alpine:3.23.2
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]