# Этап сборки: используем официальный образ Golang 1.23
FROM golang:1.23 AS builder
WORKDIR /app

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта (включая SQL-миграцию и entrypoint.sh, если нужно)
COPY . .

# Сборка приложения. Замените путь к вашему основному файлу, если он отличается.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp ./cmd/main.go

# Этап рантайма: используем минимальный образ
FROM alpine:latest
# Устанавливаем сертификаты и клиент PostgreSQL для выполнения миграций
RUN apk --no-cache add ca-certificates postgresql-client

# Устанавливаем рабочую директорию, чтобы относительные пути совпадали.
WORKDIR /app

# Копируем скомпилированное приложение из этапа сборки.
COPY --from=builder /app/myapp .

# Копируем каталог с конфигурацией (если он нужен для работы приложения).
COPY --from=builder /app/internal/config ./internal/config

# Копируем SQL-миграцию
COPY ./migrations/schema.sql ./migrations/schema.sql

# Копируем entrypoint-скрипт и делаем его исполняемым.
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Открываем порт, как в настройках приложения.
EXPOSE 8080

# Запускаем приложение через entrypoint, который сначала выполнит миграцию.
CMD ["./entrypoint.sh"]
