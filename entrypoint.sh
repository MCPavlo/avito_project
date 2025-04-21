#!/bin/sh
set -e

# Экспортируем переменную с паролем для подключения к PostgreSQL
export PGPASSWORD=$DB_PASSWORD

echo "Проверка подключения к базе данных..."

# Ждем, пока Postgres станет доступен
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "Ожидание Postgres..."
  sleep 2
done

echo "Подключение к Postgres установлено. Запуск миграций..."

# Выполняем SQL-миграцию. Флаг ON_ERROR_STOP=1 останавливает выполнение при возникновении ошибки.
psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f ./migrations/schema.sql

echo "Миграция завершена. Старт приложения..."
exec ./myapp

