# Users API

REST API in Go — GoFiber · PostgreSQL · SQLC · Uber Zap · go-playground/validator

---

## Run with Docker (recommended)

```bash
docker compose up --build
# API ready at http://localhost:8080
```

## Run locally

```bash
# 1. Install deps
go mod tidy

# 2. Create DB and run migration
psql -U postgres -c "CREATE DATABASE usersdb;"
psql -U postgres -d usersdb -f db/migrations/001_create_users_table.sql

# 3. Set env vars
cp .env.example .env   # then edit if needed

# 4. Start server
go run ./cmd/server/main.go
```

## Run tests

```bash
go test ./internal/service/... -v
```

---

## API

| Method | Path | Description |
|--------|------|-------------|
| POST | /users | Create user |
| GET | /users/:id | Get user + age |
| PUT | /users/:id | Update user |
| DELETE | /users/:id | Delete user |
| GET | /users?page=1&limit=10 | List users (paginated) |

### Examples

```bash
# Create
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'

# Get (includes age)
curl http://localhost:8080/users/1

# Update
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","dob":"1991-03-15"}'

# Delete
curl -X DELETE http://localhost:8080/users/1

# List
curl "http://localhost:8080/users?page=1&limit=10"
```
