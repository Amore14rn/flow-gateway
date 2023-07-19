FROM golang:1.20 AS build

ARG APP_PKG_NAME=flow-gateway

# Определение метаданных для образа
LABEL service="proxy-api-gateway"

# Установка рабочей директории
WORKDIR /go/src/gitlab.com/mildd/$APP_PKG_NAME

# Копирование исходного кода в контейнер
COPY . .

# Определение аргументов сборки
ARG HASHCOMMIT
ARG VERSION

# Сборка Go-приложения
RUN go build -mod=mod -v -o /out/migration ./tools/migration/*.go
RUN go build -mod=mod -v \
    -o /out/flow-gateway \
    -ldflags "-extldflags '-static' -X 'main.serviceVersion=$VERSION' -X 'main.hashCommit=$HASHCOMMIT'" \
    ./cmd/flow-gateway/*.go

RUN go mod download
COPY . .

# Сборка статического исполняемого файла без CGO
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o .

# Этап сборки завершен, создаем финальный образ
FROM alpine:latest

# Установка сертификатов CA
RUN apk add --no-cache ca-certificates

# Установка рабочей директории
WORKDIR /app

# Копирование исполняемых файлов и миграций из предыдущего образа
COPY --from=build /out/flow-gateway /app/flow-gateway
COPY --from=build /out/migration /app/migration

# Копирование миграций из исходного кода
COPY ./migrations /app/migrations

# Запуск миграций и приложения
CMD /app/migration -dir /app/migrations/sql up && /app/flow-gateway
