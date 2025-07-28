FROM golang:alpine AS builder

WORKDIR /app

# Устанавливаем инструмент для миграций
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и миграции
COPY . .
COPY migrations ./migrations

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /sas ./cmd/main.go

FROM alpine:latest
WORKDIR /app

# Копируем бинарник, миграции и конфиги
COPY --from=builder /sas .
COPY --from=builder /app/migrations ./migrations
COPY config ./config

# Устанавливаем инструмент для миграций в финальный образ
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

EXPOSE 8080

# Запускаем миграции и приложение
CMD ["sh", "-c", "migrate -path ./migrations -database \"$DB_URL\" up && migrate -path ./migrations -database \"$DB_URL\" version && ./sas"]
