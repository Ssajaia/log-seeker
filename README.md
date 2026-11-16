# LogSeeker

LogSeeker is a small log viewer with a Go HTTP server, SQLite storage, and a static HTML/CSS/JavaScript frontend. The current backend supports reading and filtering logs. It does not yet collect logs from agents or expose a log-ingestion API.

## Run with Docker Compose

Requires Docker Desktop with Compose support. From the project directory, run:

```sh
docker compose up --build
```

Open <http://localhost:8080>. Compose starts a one-time initializer before the server. If the SQLite database is empty, it inserts 30 clearly labeled demo logs; if it already contains logs, it leaves them alone. The database is stored in the `logseeker-data` named volume and persists across restarts. Keep the terminal open while using the app and press `Ctrl+C` to stop it.

To clear the database and regenerate the demo logs:

```sh
docker compose down --volumes
docker compose up --build
```

## Run locally

Requires Go 1.27 or newer. From the project directory:

```sh
go run . mock-logs
go run .
```

The first command creates 30 demo entries in `logseeker.db` if it has no logs. The second starts the server at <http://localhost:8080>. Set `LOGSEEKER_ADDR` to change the listen address or `LOGSEEKER_DB` to change the SQLite file path.

Run the existing Go tests with:

```sh
go test ./...
```

## API

`GET /api/logs` returns the newest matching entries first. Example response:

```json
{
  "logs": [
    {
      "id": 42,
      "timestamp": "2025-01-02T03:04:05Z",
      "level": "ERROR",
      "source": "worker-1",
      "message": "connection timed out"
    }
  ],
  "limit": 100,
  "offset": 0
}
```

Query parameters:

- `level`: exact severity filter; matching is case-insensitive.
- `regex`: Go/RE2 expression matched against the message before pagination. Invalid expressions return HTTP 400 with an `error` field.
- `limit`: page size, defaults to 100 and must be between 1 and 500.
- `offset`: number of matching entries to skip; defaults to 0.

Database failures return HTTP 500 with a JSON `error` field. The frontend API URL is configured in `frontend/config.js`.
