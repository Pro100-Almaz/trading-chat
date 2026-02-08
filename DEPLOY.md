# Trading Chat API — Dev Server Deployment

## Prerequisites

- Docker >= 20.10
- Docker Compose >= 2.0
- Git

## Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/Pro100-Almaz/trading-chat.git
cd trading-chat/backend

# 2. Create environment file
cp .env.example .env

# 3. Edit .env — set your secrets and config
nano .env    # or vim / any editor

# 4. Start everything
docker compose up -d --build

# 5. Verify
curl http://localhost:8080/swagger/index.html
```

The API will be available at `http://<server-ip>:8080`.

---

## Environment Configuration

Before starting, edit `.env` and update at minimum:

| Variable | What to change |
|---|---|
| `DB_HOST` | Keep `db` when using docker-compose (service name) |
| `DB_PASS` | Set a strong database password |
| `ACCESS_TOKEN_SECRET` | Random string, e.g. `openssl rand -hex 32` |
| `REFRESH_TOKEN_SECRET` | Random string, different from access secret |
| `GOOGLE_CLIENT_ID` | Your Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | Your Google OAuth client secret |
| `SMTP_*` | SMTP credentials if email verification is needed |

Generate secrets:

```bash
echo "ACCESS_TOKEN_SECRET=$(openssl rand -hex 32)" >> .env
echo "REFRESH_TOKEN_SECRET=$(openssl rand -hex 32)" >> .env
```

---

## Docker Compose Services

| Service | Container | Port | Description |
|---|---|---|---|
| `app` | `trading-chat-api` | 8080 | Go API server |
| `db` | `trading-chat-db` | 5433 (host) → 5432 (container) | PostgreSQL 17 |

---

## Common Commands

```bash
# Start in background
docker compose up -d --build

# View logs
docker compose logs -f app
docker compose logs -f db

# Stop everything
docker compose down

# Stop and remove volumes (resets database)
docker compose down -v

# Rebuild after code changes
docker compose up -d --build app

# Open psql shell
docker exec -it trading-chat-db psql -U postgres -d trading_chat

# Check service health
docker compose ps
```

---

## Running Without Docker (local development)

```bash
# Start only the database
docker compose up -d db

# In .env, set DB_HOST=localhost and DB_PORT=5433
# Then run the app natively
go run cmd/main.go
```

---

## Updating the Deployment

```bash
cd trading-chat/backend
git pull origin main
docker compose up -d --build
```

---

## Swagger API Docs

After startup, interactive API documentation is available at:

```
http://<server-ip>:8080/swagger/index.html
```

---

## Troubleshooting

**App fails to start with database connection error**
- Check that `DB_HOST=db` in `.env` (not `localhost`) when using docker-compose
- Wait for db health check: `docker compose ps` should show db as `healthy`

**Port 8080 already in use**
- Change `PORT` in `.env` and update `ports` mapping accordingly
- Or stop the conflicting process: `lsof -i :8080`

**Database data reset**
- `docker compose down -v` removes the pgdata volume
- Next `docker compose up` creates a fresh database (migrations run on app start)

**View container resource usage**
```bash
docker stats trading-chat-api trading-chat-db
```
