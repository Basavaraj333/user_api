

## Run with Docker (Recommended)

```bash
# Clone the repo
git clone https://github.com/YOUR_USERNAME/users-api.git
cd users-api

# Start everything — PostgreSQL + Go server
docker compose up --build
```

API is live at `http://localhost:8080`



## Run Locally (without Docker)

```bash
# 1. Install dependencies
go mod tidy

# 2. Create database
psql -U postgres -c "CREATE DATABASE usersdb;"

# 3. Run migration
psql -U postgres -d usersdb -f db/migrations/001_create_users_table.sql

# 4. Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=usersdb
export SERVER_PORT=8080
export APP_ENV=development

# 5. Run server
go run ./cmd/server/main.go
```



## Run Tests

```bash
go test ./internal/service/... -v
```


## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /users | Create a new user |
| GET | /users/:id | Get user by ID (includes age) |
| PUT | /users/:id | Update user |
| DELETE | /users/:id | Delete user (204 No Content) |
| GET | /users?page=1&limit=10 | List all users with pagination |



## API Examples

### Create User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```
Response `201`:
```json
{"id": 1, "name": "Alice", "dob": "1990-05-10"}
```

### Get User with Age
```bash
curl http://localhost:8080/users/1
```
Response `200`:
```json
{"id": 1, "name": "Alice", "dob": "1990-05-10", "age": 36}
```

### Update User
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","dob":"1991-03-15"}'
```
Response `200`:
```json
{"id": 1, "name": "Alice Updated", "dob": "1991-03-15"}
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/users/1
```
Response: `204 No Content`

### List Users
```bash
curl "http://localhost:8080/users?page=1&limit=10"
```
Response `200`:
```json
{
  "data": [{"id": 1, "name": "Alice", "dob": "1990-05-10", "age": 36}],
  "total": 1,
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```



## Design Decisions

| Decision | Reason |
|----------|--------|
| Layered architecture | Each layer has one job — easy to test and maintain |
| SQLC over GORM | Raw SQL is explicit and type-safe — no ORM magic |
| Interface for repo and service | Enables mock injection for unit testing |
| Age calculated at runtime | DOB is a fact, age is derived — no need to store derived data |
| Sentinel error ErrNotFound | Clean separation between DB errors and domain errors |
| Multi-stage Dockerfile | Final image is tiny — only the binary, no Go toolchain |
| UUID requestId per request | Enables tracing a single request across all log lines |