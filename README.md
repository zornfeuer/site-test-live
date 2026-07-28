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

Триггеры: каждый push в `master`, а также ручной запуск (`workflow_dispatch`).
Образ всегда публикуется под тегом `latest`.

Аутентификация в GHCR при публикации — встроенный `GITHUB_TOKEN`, дополнительных секретов для
самой сборки настраивать не нужно.

Локально запустить сборку так же, как в пайплайне:

```bash
docker build -t ghcr.io/<owner>/<repo>-backend:local ./backend
docker build -t ghcr.io/<owner>/<repo>-frontend:local ./frontend
```

### Деплой на VPS

После публикации образов job `deploy` заходит по SSH на VPS и накатывает свежую версию:

```bash
cd /root/site-test-live
git pull --ff-only origin master
docker login ghcr.io -u <actor> --password-stdin   # токеном GITHUB_TOKEN текущего запуска
docker compose pull
docker compose up -d
docker image prune -f
```

Предпосылки на самой VPS (настраиваются один раз вручную, вне GitHub Actions):

- репозиторий уже склонирован в `/root/site-test-live`, у `git pull` есть доступ до этого приватного
  репозитория (например, через deploy key, добавленный в GitHub — Settings → Deploy keys);
- установлены Docker и плагин Docker Compose (`docker compose version`);
- публичный ключ, соответствующий секрету `VPS_SSH_KEY` (см. ниже), добавлен в
  `~/.ssh/authorized_keys` пользователя, под которым идёт деплой.

Пакеты в GHCR по умолчанию **приватные** и привязаны к репозиторию — `docker login` на VPS решает эту
проблему через `GITHUB_TOKEN`, дополнительный Personal Access Token не нужен. Если хочется публичный
анонимный `docker pull`, можно вместо этого открыть Settings пакета на GitHub и переключить видимость
на Public — тогда шаг логина на VPS можно будет убрать.

#### Переменные окружения (GitHub → Settings → Secrets and variables → Actions → Secrets)

| Секрет | Обязателен | По умолчанию | Назначение |
|---|---|---|---|
| `VPS_HOST` | да | — | IP-адрес или домен VPS |
| `VPS_SSH_KEY` | да | — | Приватный SSH-ключ целиком (включая строки `-----BEGIN ... KEY-----` / `-----END ... KEY-----`), которым пайплайн заходит на VPS |
| `VPS_USER` | нет | `root` | Пользователь для SSH-подключения |
| `VPS_PORT` | нет | `22` | SSH-порт |

`GITHUB_TOKEN` создаётся автоматически для каждого запуска и отдельно настраивать его не нужно —
он используется и для публикации образов, и для `docker login` на VPS при деплое.
