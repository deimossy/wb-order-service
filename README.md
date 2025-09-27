# WB Order Service

Демонстрационный микросервис для работы с заказами на Go с использованием PostgreSQL, Kafka и кеша в памяти.

---

## Описание

Сервис получает данные о заказах из Kafka, сохраняет их в PostgreSQL и кэширует в памяти для быстрого доступа. Предоставляется HTTP API и веб-интерфейс для просмотра заказов по `order_uid`.

Возможности:

* Подписка на Kafka и обработка сообщений о заказах.
* Сохранение заказов в PostgreSQL.
* Кэширование последних заказов в памяти.
* Восстановление кеша из БД при старте.
* HTTP API для получения заказа в JSON.
* Веб-интерфейс для запроса заказа по ID.

---

## Технологии

* Go
* PostgreSQL
* Kafka
* HTML/JS
* Docker/Docker Compose
* Makefile

---

## Установка и запуск

### 1. Клонирование репозитория

```bash
git clone https://github.com/deimossy/wb-order-service.git
cd wb-order-service
```

### 2. Настройка `.env`

Создайте `.env` на основе `.env.example` и укажите параметры подключения к PostgreSQL и Kafka.

### 3. Сборка и запуск

1. **Сборка Docker-образов:**

```bash
make build
```

2. **Запуск Kafka и приложения через Docker Compose:**

```bash
make up
```

3. **Отправка тестовых заказов из JSON-файлов:**

```bash
make send-orders
```

Сервис будет доступен по адресу `http://localhost:8080`.

### 4. Остановка сервиса

```bash
make down
```

---

## Makefile — цели и пояснения

| Цель                | Описание                                                                                      |
| ------------------- | --------------------------------------------------------------------------------------------- |
| `test`              | Запуск всех тестов Go с подробным выводом                                                     |
| `bench`             | Запуск benchmark тестов для `order_service`                                                   |
| `create-network`    | Создание Docker-сети `infra-net`                                                              |
| `build`             | Сборка Docker-образов приложения и Kafka                                                      |
| `up`                | Запуск Kafka (инициализация в корневой папке `kafka/`) и приложения через Docker Compose      |
| `down`              | Остановка и удаление контейнеров Kafka и приложения                                           |
| `send-orders`       | Отправка JSON-файлов из `internal/services/order_service/testdata` в Kafka топик `get_orders` |
| `logs-app-docker`   | Просмотр логов приложения в Docker Compose                                                    |
| `logs-kafka-docker` | Просмотр логов Kafka в Docker Compose                                                         |
| `logger`            | Просмотр логов контейнера `order_service` напрямую                                            |
| `clean`             | Полная очистка Docker Compose и удаление сети `infra-net`                                     |

> **Примечание:** Для `send-orders` можно добавлять новые JSON-файлы в папку `testdata` сервиса заказов. Makefile автоматически пройдет по всем JSON и отправит их в Kafka.

---

## Работа с сервисом

### HTTP API

Получение заказа:

```
GET http://localhost:8080/order/<order_uid>
```

Возвращает JSON с данными заказа.

### Веб-интерфейс

* Откройте `http://localhost:8080`.
* Введите `order_uid` и получите информацию о заказе.

---

## Структура проекта

* `cmd/` — точка входа в приложение.
* `internal/` — внутренние пакеты:

  * `kafka/` — консюмер Kafka для обработки заказов.
  * `services/` — логика работы с заказами, кешем и БД.
  * `services/order_service/testdata/` — папка для JSON-файлов заказов, используемых в `send-orders`.
* `migrations/` — миграции базы данных.
* `pkg/` — вспомогательные библиотеки:

  * Кастомные ошибки.
  * Клиент для работы с БД.
  * Логгер для приложения.
* `ui/` — веб-интерфейс.
* `kafka/` (корневая папка) — файлы для инициализации Kafka (docker-compose, настройки).
* `Dockerfile`, `docker-compose.yaml`, `Makefile` — сборка и запуск проекта.

---

## Модель данных заказа

```json
{
  "order_uid": "123e4567-e89b-12d3-a456-426655440000",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```

---

## Примечания

* Некорректные сообщения из Kafka логируются и игнорируются.
* Данные не теряются при сбоях благодаря транзакциям и подтверждению сообщений.
* Кеш ускоряет повторные запросы.
* Последовательность запуска: **сначала `make build`, затем `make up`, после этого `make send-orders`**.

---

## Лицензия

MIT
