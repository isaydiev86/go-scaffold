# go-scaffold

[![GitHub release](https://img.shields.io/github/v/release/isaydiev86/go-scaffold)](https://github.com/isaydiev86/go-scaffold/releases)
[![Go version](https://img.shields.io/badge/go-1.26+-blue)](https://golang.org)

CLI-инструмент для генерации production-ready Go микросервисов по шаблону Clean Architecture.

## Установка

```bash
# последняя версия
go install github.com/isaydiev86/go-scaffold/cmd/go-scaffold@latest

# конкретная версия
go install github.com/isaydiev86/go-scaffold/cmd/go-scaffold@v1.0.0
```

Или собрать из исходников:

```bash
git clone https://github.com/isaydiev86/go-scaffold
cd go-scaffold
go build -o go-scaffold ./cmd/go-scaffold
```

## Использование

### Интерактивный режим

Запустите команду без аргументов — инструмент задаст вопросы:

```
$ go-scaffold new

? Service name: payments
? Go module path [github.com/myorg/payments]: github.com/mycompany/payments
? Output directory [payments]:
? Add Redis? (caching example for GET) y/n [n]: y
? Add Kafka consumer? y/n [n]: n
? Add Kafka producer? (transactional outbox + worker) y/n [n]: y

Generating service "payments" → payments

Done! Next steps:
  cd payments
  make bootstrap   # go mod tidy + go generate (mocks, OpenAPI) — required before the first build
  go test ./...
  make up          # start postgres + redis + kafka
  make migrate-up  # apply migrations
  make run         # start the service
```

> Сгенерированный проект собирается только после кодогенерации (`make bootstrap`): oapi-codegen создаёт HTTP-сервер по OpenAPI-спецификации, mockery — моки. Для этого нужны установленные инструменты: `make oapi-install mockery-install` (однократно).

### Передача аргументов напрямую

```bash
# только имя сервиса — модуль и директорию спросит интерактивно
go-scaffold new payments

# полностью без вопросов (удобно для скриптов и CI)
go-scaffold new payments \
  --module github.com/mycompany/payments \
  --output ./services/payments \
  --redis \
  --kafka-producer
```

### CI-режим (`--non-interactive`)

С флагом `--non-interactive` команда никогда не читает stdin: отсутствующие обязательные значения (имя сервиса, `--module`) дают ошибку, boolean-флаги без явного указания остаются выключенными. Пример шага CI:

```bash
go-scaffold new payments \
  --module github.com/mycompany/payments \
  --output ./services/payments \
  --non-interactive

cd ./services/payments
make verify   # bootstrap (go mod tidy + go generate + go mod tidy) + go test ./...
```

### Флаги команды `new`

| Флаг | Описание | По умолчанию |
|------|----------|--------------|
| `--module` | Go module path | `github.com/myorg/<service-name>` |
| `--output` | Директория для генерации | имя сервиса |
| `--redis` | Подключить Redis (кэширование GET-запроса) | без Redis (спросит интерактивно) |
| `--kafka-consumer` | Подключить Kafka consumer (контроллер с consumer group) | без консюмера (спросит интерактивно) |
| `--kafka-producer` | Подключить Kafka producer + транзакционный outbox + воркер доставки | без продьюсера (спросит интерактивно) |
| `--non-interactive` | Никогда не задавать вопросов: ошибка при отсутствии обязательных значений (для CI) | интерактивный режим |

Имя сервиса валидируется: строчные латинские буквы, цифры и одиночные дефисы; начинается с буквы, не заканчивается дефисом (например, `payments`, `new-payments`).

Опции независимы: можно сгенерировать сервис только с консюмером, только с продьюсером или с обоими. Producer всегда идёт в связке с outbox-паттерном: usecase `Create` пишет событие в таблицу `outbox` в той же транзакции, что и сущность, а фоновый воркер батчами доставляет события в Kafka (at-least-once).

## Структура генерируемого проекта

```
payments/
├── cmd/app/                        # Точка входа приложения
├── config/                         # Загрузка конфигурации из .env
├── internal/
│   ├── domain/                     # Сущности, статусы, ошибки (+ Event с --kafka-producer)
│   ├── usecase/                    # Бизнес-логика, интерфейсы адаптеров
│   ├── adapter/
│   │   ├── postgres/               # Реализация репозитория (pgx)
│   │   ├── redis/                  # Кэш (только с --redis)
│   │   └── kafka_producer/         # Продьюсер kafka-go (только с --kafka-producer)
│   ├── controller/
│   │   ├── http/v1/                # HTTP-обработчики (Chi)
│   │   ├── kafka_consumer/         # Консюмер kafka-go (только с --kafka-consumer)
│   │   └── worker/                 # Outbox-воркер (только с --kafka-producer)
│   └── dto/                        # Data Transfer Objects
├── pkg/                            # Переиспользуемые пакеты (logger, postgres, ...)
├── api/http/                       # OpenAPI спецификация
├── gen/                            # Сгенерированный код (oapi-codegen)
├── migration/postgres/             # SQL-миграции (+ таблица outbox с --kafka-producer)
├── Makefile
├── Dockerfile
├── docker-compose.yml              # PostgreSQL (+ Redis с --redis; + Kafka и Kafka UI с --kafka-*)
└── .env.example
```

## Команды в сгенерированном проекте

```bash
make bootstrap        # go mod tidy + go generate + go mod tidy (обязательно после генерации)
make verify           # bootstrap + go test ./...
make run              # запустить сервис
make test             # тесты с покрытием
make integration-test # интеграционные тесты
make lint             # golangci-lint
make generate         # go generate (моки, OpenAPI)
make up               # docker-compose up (postgres; с --redis также redis; с --kafka-* также kafka)
make migrate-up       # применить миграции
make migrate-down     # откатить миграции
make migrate-create   # создать новую миграцию
```

## Технологии в шаблоне

| Категория | Библиотека |
|-----------|------------|
| HTTP роутер | [go-chi/chi](https://github.com/go-chi/chi) |
| PostgreSQL | [jackc/pgx](https://github.com/jackc/pgx) |
| Redis (опционально) | [redis/go-redis](https://github.com/redis/go-redis) |
| Kafka (опционально) | [segmentio/kafka-go](https://github.com/segmentio/kafka-go) |
| Миграции | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Логирование | [rs/zerolog](https://github.com/rs/zerolog) |
| Валидация | [go-playground/validator](https://github.com/go-playground/validator) |
| Конфиг | [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig) |
| Моки | [vektra/mockery](https://github.com/vektra/mockery) |
| OpenAPI | [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) |

## Требования

- Go 1.26+
