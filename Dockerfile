# Этап сборки: используем официальный образ Golang 1.23
FROM golang:1.23 AS builder
WORKDIR /app

# Переносим файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Сборка приложения. Замените путь к вашему основному файлу, если он отличается.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp ./cmd/main.go

# Этап рантайма: используем минимальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию, например /app, чтобы относительные пути совпадали.
WORKDIR /app

# Копируем скомпилированное приложение
COPY --from=builder /app/myapp .

# Копируем каталог с конфигурацией
COPY --from=builder /app/internal/config ./internal/config

# Открываем порт, как и было настроено
EXPOSE 8080

# Запускаем приложение
CMD ["./myapp"]
