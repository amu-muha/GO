# Workout API

A Go HTTP API for creating, reading, updating, and deleting workouts and their exercise entries. The API uses PostgreSQL for persistence and applies embedded Goose migrations automatically when it starts.

## Tech Stack

- Go 1.26+
- PostgreSQL 12+
- Chi router
- Goose database migrations
- pgx PostgreSQL driver

## Prerequisites

- Go 1.26 or newer
- Docker and Docker Compose

## Getting Started

1. Start the development PostgreSQL database:

   ```bash
   docker compose up -d db
   ```

   The database is available at `localhost:5434` with these credentials:

   - Database: `postgres`
   - User: `postgres`
   - Password: `postgresss`

2. Configure the database connection. The repository's local `.env` uses:

   ```bash
   export DATABASE_URL='host=localhost user=postgres password=postgresss dbname=postgres port=5434 sslmode=disable'
   ```

   The application reads `DATABASE_URL` from the environment; it does not load `.env` automatically.

3. Download dependencies and start the API:

   ```bash
   go mod download
   go run .
   ```

   The server listens on `http://localhost:8080` by default. Use a different port with:

   ```bash
   go run . -port 8081
   ```

Migrations in `migration/` run automatically during application startup.

## API

All JSON responses use an envelope such as `{ "workout": ... }`.

### Health check

```http
GET /health
```

### Create a workout

```http
POST /workouts
Content-Type: application/json
```

Example request:

```json
{
  "title": "Full body strength",
  "description": "A balanced strength session",
  "duration_minutes": 45,
  "calories_burned": 320,
  "entries": [
    {
      "exercise_name": "Back squat",
      "sets": 4,
      "reps": 8,
      "notes": "Keep the tempo controlled",
      "order_index": 1
    },
    {
      "exercise_name": "Plank",
      "sets": 3,
      "duration_seconds": 45,
      "order_index": 2
    }
  ]
}
```

Each workout entry must provide exactly one of `reps` or `duration_seconds`. `weight` is optional.

### Get a workout

```http
GET /workouts/{id}
```

### Update a workout

```http
PUT /workouts/{id}
Content-Type: application/json
```

The update body may include `title`, `description`, `duration_minutes`, `calories_burned`, and `entries`. Sending `entries` replaces the workout's existing entries.

### Delete a workout

```http
DELETE /workouts/{id}
```

Returns `204 No Content` when the workout is deleted.

## Example Commands

Check that the server is running:

```bash
curl http://localhost:8080/health
```

Create a workout:

```bash
curl -X POST http://localhost:8080/workouts \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Morning run",
    "description": "Easy-paced run",
    "duration_minutes": 30,
    "calories_burned": 250,
    "entries": [{
      "exercise_name": "Running",
      "sets": 1,
      "duration_seconds": 1800,
      "order_index": 1
    }]
  }'
```

Retrieve workout `1`:

```bash
curl http://localhost:8080/workouts/1
```

## Database Services

`docker-compose.yml` defines two PostgreSQL services:

- `db`: development database on host port `5434`
- `test_db`: test database on host port `5433`

Start both services with:

```bash
docker compose up -d
```

Stop the containers with:

```bash
docker compose down
```

The PostgreSQL data directories are mounted under `database/` and are ignored by Git.

## Project Structure

```text
.
├── main.go                    # HTTP server entry point
├── internal/
│   ├── api/                   # HTTP handlers
│   ├── app/                   # Application initialization
│   ├── routes/                # Route registration
│   ├── store/                 # PostgreSQL access and workout models
│   └── utils/                 # JSON and URL parameter helpers
├── migration/                 # Embedded Goose SQL migrations
├── docker-compose.yml         # Development and test databases
└── go.mod                     # Go module and dependencies
```

## Development Checks

Run the package tests and build the API with:

```bash
go test ./...
go build ./...
```