# Notes demo

Тестовое демонстрационное приложение: список заметок.

- **Frontend**: Vue 3 + TypeScript (Vite)
- **Backend**: Go (`net/http` + `pgx`)
- **DB**: PostgreSQL

## Запуск через Docker Compose (тестовый стенд)

```bash
docker compose up --build
```

- Frontend: http://localhost:8081
- Backend API: http://localhost:8090/api/notes
- Postgres: localhost:5433 (demo/demo, база `demo`)

> Порты подобраны так, чтобы не конфликтовать с локально уже запущенными сервисами.
> При необходимости поменяйте их в `docker-compose.yml`.

Схема БД (`backend/schema.sql`) применяется автоматически при первом старте контейнера `db`.

Остановить и удалить контейнеры (данные БД сохранятся в volume `db-data`):

```bash
docker compose down
```

Удалить вместе с данными:

```bash
docker compose down -v
```

## Локальная разработка без Docker

Поднять только БД в Docker, а backend/frontend запускать напрямую:

```bash
docker compose up db
```

### Backend

```bash
cd backend
DATABASE_URL="postgres://demo:demo@localhost:5433/demo?sslmode=disable" go run .
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Откроется на http://localhost:5173, запросы к `/api/*` проксируются на backend (см. `vite.config.ts`).

## API

| Метод  | Путь              | Описание              |
|--------|-------------------|------------------------|
| GET    | `/api/notes`      | список заметок         |
| POST   | `/api/notes`      | создать `{ "text": "..." }` |
| DELETE | `/api/notes/{id}` | удалить заметку        |
| GET    | `/api/health`     | health-check           |
