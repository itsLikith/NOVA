# Auth Service

A Go-based authentication service for NOVA, built with Fiber and PostgreSQL. It handles user authentication, JWT issuance/validation, admin onboarding, and protected route access checks.

## Overview

This service is responsible for:

- authenticating users with username and password
- issuing signed JWTs
- validating bearer tokens for protected endpoints
- creating regular users through admin-only access
- auto-seeding the initial admin user from environment variables
- auto-creating the PostgreSQL database if it does not already exist

## Tech Stack

- Go 1.26+
- Fiber v2
- PostgreSQL
- GORM
- JWT (golang-jwt/jwt/v5)
- bcrypt

## Project Structure

```text
src/server/services/auth/
├── cmd/
│   └── server/
│       └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── handler/
│   │   └── handler.go
│   ├── model/
│   │   └── model.go
│   ├── repository/
│   │   ├── postgresql.go
│   │   └── repository.go
│   └── service/
│       └── service.go
├── pkg/
│   ├── database/
│   │   └── database.go
│   ├── middleware/
│   │   └── auth.go
│   └── utils/
│       └── utils.go
├── tests/
│   └── service_test.go
├── go.mod
├── README.md
└── .env.example (recommended to add at project level or service level)
```

## Configuration

The service reads configuration from environment variables instead of hardcoded values.

### Required environment variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=nova_auth

APP_PORT=8081
JWT_SECRET=super_secure_jwt_secret

ADMIN_USERNAME=admin
ADMIN_PASSWORD=strong_admin_password
```

### Default behavior

If a value is not provided, the service uses local development defaults for database access, but admin credentials and JWT secret should always be configured in real deployments.

## Database Setup

On startup, the service:

1. builds the PostgreSQL DSN from config
2. checks whether the target database exists
3. creates it automatically if needed
4. runs GORM auto-migration for the `User` model

This reduces startup friction for local and containerized development.

## Admin Bootstrap

The service seeds an initial admin user at startup using `ADMIN_USERNAME` and `ADMIN_PASSWORD` from the environment.

If the admin user already exists, the service skips creation. If it does not exist, it hashes the configured password with bcrypt and creates the admin record with the `admin` role.

## Authentication Flow

### Login

`POST /api/v1/auth/login`

Request body:

```json
{
  "username": "admin",
  "password": "strong_admin_password"
}
```

Response:

```json
{
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_at": "2026-08-16T12:00:00Z",
    "user": {
      "id": 1,
      "username": "admin",
      "role": "admin",
      "created_at": "2026-08-16T00:00:00Z"
    }
  }
}
```

### Validate token

`GET /api/v1/auth/validate`

Headers:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "message": "token is valid",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "role": "admin"
    },
    "expires_at": "2026-08-16T12:00:00Z"
  }
}
```

### Create user

`POST /api/v1/auth/users`

This route is protected and requires admin access.

Headers:

```http
Authorization: Bearer <admin_token>
```

Request body:

```json
{
  "username": "alice",
  "password": "secret123"
}
```

The admin must send a valid bearer token in the Authorization header and the token must belong to an admin user.

Response:

```json
{
  "message": "user created successfully",
  "data": {
    "id": 2,
    "username": "alice",
    "role": "user",
    "created_at": "2026-08-16T00:00:00Z"
  }
}
```

Error responses use the same envelope:

```json
{
  "message": "authentication required",
  "error": "missing authorization header"
}
```

## Middleware

The service includes middleware for route protection:

- `AuthRequired(secret)` checks JWT validity and populates request locals
- `AdminRequired(secret)` ensures the token belongs to an admin user

## API Health Check

`GET /health`

Response:

```json
{
  "status": "ok",
  "service": "auth"
}
```

## Local Development

From the service directory:

```bash
cd src/server/services/auth

go mod tidy
go run cmd/server/main.go
```

This will start the auth service on the configured port (default is `8081`).

## Docker Compose

The project root includes a Postgres container definition. Example:

```yaml
services:
  auth_db:
    image: postgres:18.4-alpine3.24
    environment:
      POSTGRES_USER: ${DB_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-postgres}
      POSTGRES_DB: ${DB_NAME:-nova_auth}
    ports:
      - "5432:5432"
```

Run:

```bash
docker compose up -d auth_db
```

## Running Tests

The service tests are isolated under the `tests` directory.

```bash
cd src/server/services/auth
go test ./...
```

## Security Notes

- JWT secret should be stored in a secure environment manager or secret store
- admin credentials should never be committed to source control
- use strong passwords in real environments
- treat the auth service as a trust boundary and protect all routes behind proper middleware

## Notes

The service is intentionally simple and modular, making it easy to expand for:

- email/password registration flows
- refresh tokens
- password reset workflows
- role-based authorization expansion
- API gateway integration

## Typical Startup Flow

When the app starts:

1. load env config
2. validate required values
3. connect to PostgreSQL
4. create the database if missing
5. auto-migrate `User`
6. seed admin user from env
7. start Fiber server on configured port

This makes the service reliable for both local development and staged deployment environments.
