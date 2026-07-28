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

## CI/CD

Workflow [`.github/workflows/docker-publish.yml`](.github/workflows/docker-publish.yml) собирает образы
`backend` и `frontend` и публикует их в GitHub Container Registry:

- `ghcr.io/<owner>/<repo>-backend`
- `ghcr.io/<owner>/<repo>-frontend`

Триггеры: push в `main`, теги `v*.*.*`, а также ручной запуск (`workflow_dispatch`).
Тегирование образов (через `docker/metadata-action`): имя ветки, semver-тег (для тегов вида `v1.2.3`),
короткий SHA коммита и `latest` — только для `main`.

Аутентификация — встроенный `GITHUB_TOKEN`, дополнительных секретов настраивать не нужно.
После первой публикации пакет в GHCR по умолчанию **приватный** и привязан к репозиторию — если нужен
анонимный `docker pull`, откройте Settings пакета на GitHub и переключите видимость на Public.

Локально запустить сборку так же, как в пайплайне:

```bash
docker build -t ghcr.io/<owner>/<repo>-backend:local ./backend
docker build -t ghcr.io/<owner>/<repo>-frontend:local ./frontend
```
