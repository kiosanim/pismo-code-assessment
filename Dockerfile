FROM golang:1.25.6-alpine3.23 as builder
WORKDIR /app
COPY . .
RUN go build -o app_api ./cmd/api && go build -o app_migration_tool ./cmd/migrate

FROM alpine:3.23.2
WORKDIR /app
RUN mkdir -p internal/infra/database/migrations
COPY --from=builder /app/app_* .
COPY --from=builder /app/internal/infra/database/migrations/* /app/internal/infra/database/migrations
EXPOSE 8080
CMD ["./app_api"]