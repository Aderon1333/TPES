# Используем официальный образ Go в качестве базового
FROM golang:1.22.4-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum в рабочую директорию
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код в рабочую директорию
COPY . /app

# путь к main.go
WORKDIR /app/cmd/tpes
# билд приложения 
RUN go build -o /app/bin/tpes

# Определяем порт, на котором будет работать приложение
EXPOSE 8080

# Команда для запуска приложения
ENTRYPOINT ["/app/bin/tpes"]