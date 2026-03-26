# go-scaffold

CLI-инструмент для генерации production-ready Go микросервисов по шаблону Clean Architecture.

## Установка

```bash
go install github.com/isaydiev86/go-scaffold/cmd/go-scaffold@latest
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

Generating service "payments" → payments

Done! Next steps:
  cd payments
  go mod tidy
  make up          # start postgres + redis
  make migrate-up  # apply migrations
  make run         # start the service
```

### Передача аргументов напрямую

```bash
# только имя сервиса — модуль и директорию спросит интерактивно
go-scaffold new payments

# полностью без вопросов (удобно для скриптов и CI)
go-scaffold new payments \
  --module github.com/mycompany/payments \
  --output ./services/payments
```

### Флаги команды `new`

| Флаг | Описание | По умолчанию |
|------|----------|--------------|
| `--module` | Go module path | `github.com/myorg/<service-name>` |
| `--output` | Директория для генерации | имя сервиса |

## Структура генерируемого проекта

```
payments/
├── cmd/app/                        # Точка входа приложения
├── config/                         # Загрузка конфигурации из .env
├── internal/
│   ├── domain/                     # Сущности, статусы, ошибки
│   ├── usecase/                    # Бизнес-логика, интерфейсы адаптеров
│   ├── adapter/
│   │   ├── postgres/               # Реализация репозитория (pgx)
│   │   └── redis/                  # Кэш и идемпотентность
│   ├── controller/http/v1/         # HTTP-обработчики (Chi)
│   └── dto/                        # Data Transfer Objects
├── pkg/                            # Переиспользуемые пакеты (logger, postgres, redis, ...)
├── api/http/                       # OpenAPI спецификация
├── gen/                            # Сгенерированный код (oapi-codegen)
├── migration/postgres/             # SQL-миграции (golang-migrate)
├── Makefile
├── Dockerfile
├── docker-compose.yml              # PostgreSQL, Redis, Redis Commander
└── .env.example
```

## Команды в сгенерированном проекте

```bash
make run              # запустить сервис
make test             # тесты с покрытием
make integration-test # интеграционные тесты
make lint             # golangci-lint
make generate         # go generate (моки, OpenAPI)
make up               # docker-compose up (postgres + redis)
make migrate-up       # применить миграции
make migrate-down     # откатить миграции
make migrate-create   # создать новую миграцию
```

## Технологии в шаблоне

| Категория | Библиотека |
|-----------|------------|
| HTTP роутер | [go-chi/chi](https://github.com/go-chi/chi) |
| PostgreSQL | [jackc/pgx](https://github.com/jackc/pgx) |
| Redis | [redis/go-redis](https://github.com/redis/go-redis) |
| Миграции | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Логирование | [rs/zerolog](https://github.com/rs/zerolog) |
| Валидация | [go-playground/validator](https://github.com/go-playground/validator) |
| Конфиг | [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig) |
| Моки | [vektra/mockery](https://github.com/vektra/mockery) |
| OpenAPI | [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) |

## Требования

- Go 1.26+
