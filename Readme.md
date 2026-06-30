# Subscriptions API
 
REST-сервис для агрегации данных об онлайн-подписках пользователей: создание, чтение, обновление и удаление записей о подписках, а также подсчёт суммарной стоимости за период.
 
Проект построен по слоистой архитектуре (`config → database → repository → handlers → server`) с ручным внедрением зависимостей. В качестве хранилища используется PostgreSQL, документация по API генерируется через Swagger.
 
## Стек
 
- **Go** — язык реализации
- **chi/v5** — HTTP-роутер и middleware (логирование, recover)
- **pgx/v5 (pgxpool)** — драйвер и пул соединений к PostgreSQL
- **slog** — структурное логирование
- **golang-migrate** — миграции схемы БД
- **swaggo / http-swagger** — генерация и отдача Swagger UI
- **Docker / Docker Compose** — контейнеризация
## Возможности
 
| Метод    | Путь                     | Описание                                   |
| -------- | ------------------------ | ------------------------------------------ |
| `POST`   | `/subscriptions`         | Создать подписку                           |
| `GET`    | `/subscriptions`         | Список всех подписок                        |
| `GET`    | `/subscriptions/{id}`    | Получить подписку по ID                     |
| `PUT`    | `/subscriptions/{id}`    | Обновить подписку                           |
| `DELETE` | `/subscriptions/{id}`    | Удалить подписку                            |
| `GET`    | `/subscriptions/total`   | Суммарная стоимость подписок за период      |
| `GET`    | `/swagger/index.html`    | Интерактивная Swagger-документация          |
 
## Быстрый старт (Docker Compose)
 
Самый простой способ — поднять приложение вместе с базой одной командой.
 
> **Важно:** перед сборкой нужно сгенерировать Swagger-документацию, так как пакет `docs` подключается в коде, но не хранится в репозитории (он в `.gitignore`). Без этого шага сборка упадёт.
 
```bash
# 1. Установить генератор Swagger (один раз)
go install github.com/swaggo/swag/cmd/swag@latest
 
# 2. Сгенерировать документацию (создаст папку docs/)
swag init -g cmd/api/main.go
 
# 3. Запустить приложение и БД
docker compose up --build
```
 
После запуска сервис доступен на `http://localhost:8080`, Swagger UI — на `http://localhost:8080/swagger/index.html`.
 
## Локальный запуск (без контейнера приложения)
 
1. Поднять только базу данных:
```bash
   docker compose up -d postgres
```
 
2. Создать файл `.env` в корне проекта (см. `env.exp`):
```env
   SQL_CONN=postgres://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable
   PORT=8080
```
 
3. Применить миграции (нужен установленный [golang-migrate](https://github.com/golang-migrate/migrate)):
```bash
   make migrate-up
```
 
4. Сгенерировать Swagger и запустить сервис:
```bash
   swag init -g cmd/api/main.go
   go run ./cmd/api
```
 
## Переменные окружения
 
| Переменная | Обязательна | По умолчанию | Описание                          |
| ---------- | ----------- | ------------ | --------------------------------- |
| `SQL_CONN` | да          | —            | строка подключения к PostgreSQL   |
| `PORT`     | нет         | `8080`       | порт HTTP-сервера                 |
 
## Миграции
 
Миграции лежат в `migrations/` и применяются через `golang-migrate`:
 
```bash
make migrate-up     # применить
make migrate-down   # откатить
```
 
Схема таблицы `subscriptions`: `id`, `service_name`, `price`, `user_id (UUID)`, `start_date`, `end_date (nullable)`.
 
## Формат даты
 
Даты передаются в формате **`MM-YYYY`** (например, `07-2025`). Это касается полей `start_date`, `end_date`, а также параметров `from` и `to` в эндпоинте подсчёта стоимости.
 
## Примеры запросов
 
Создать подписку:
 
```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```
 
Получить список подписок:
 
```bash
curl http://localhost:8080/subscriptions
```
 
Посчитать суммарную стоимость за период (с опциональным фильтром по сервису):
 
```bash
curl "http://localhost:8080/subscriptions/total?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&from=07-2025&to=12-2025&service_name=Yandex%20Plus"
```
 
## Структура проекта
 
```
.
├── cmd/api/            # точка входа (main)
├── internal/
│   ├── config/         # загрузка конфигурации из окружения
│   ├── database/       # подключение к PostgreSQL (pgxpool)
│   ├── models/         # модели данных и входные DTO
│   ├── repository/     # работа с БД (CRUD + агрегация)
│   ├── handlers/       # HTTP-обработчики
│   └── server/         # роутинг и запуск HTTP-сервера
├── migrations/         # SQL-миграции (up/down)
├── docker-compose.yml  # PostgreSQL + приложение
├── Dockerfile          # многоэтапная сборка образа
└── Makefile            # команды миграций
```
 
## Заметки по разработке
 
- Документация Swagger перегенерируется командой `swag init -g cmd/api/main.go` после изменения аннотаций в обработчиках.
- Все запросы к БД параметризованы — защита от SQL-инъекций.
- Логирование — структурное (`slog`), middleware `Recoverer` не даёт упасть серверу при панике в обработчике.
 
