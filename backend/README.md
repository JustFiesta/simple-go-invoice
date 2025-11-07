# Backend

RESTful API built with Go and Gin framework for managing invoices and generating PDFs.

## What It Does

A lightweight HTTP API that handles:

- Invoice CRUD
- Pagination & filtering - search, sort, paginate invoice lists
- PDF generation
- Validation
- HATEOAS - REST responses with hypermedia links

## Technologies

- Go 1.25
- Gin
- GORM
- SQLite
- gofpdf

## Architecture

```shell
backend/
├── main.go              # Entry point, server setup
├── config/              # Configuration loading
│   └── config.go        # Port and settings
├── database/            # Database layer
│   └── connection.go    # SQLite connection & migrations
├── models/              # Data structures
│   ├── invoice.go       # Invoice entity
│   ├── invoice_item.go  # Line item entity
│   └── response.go      # API response wrappers
├── controllers/         # Request handlers
│   ├── invoice_controller.go
│   └── invoice_items.go
├── services/            # Business logic
│   └── pdf_service.go   # PDF generation
└── routes/              # Route definitions
    └── routes.go        # Endpoint mappings + middleware
```

## Design Patterns

### MVC-like Structure

- Routes - map URLs to controllers
- Controllers - handle HTTP requests/responses
- Models - define data structures and validation
- Services - encapsulate business logic (PDFs)

### Response Format

All responses follow a consistent structure:

```json
{
  "data": {...},
  "links": {
    "self": "/api/invoices/1",
    "related": {
      "items": "/api/invoices/1/items"
    }
  },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 42
  }
}
```

Errors use the same envelope:

```json
{
  "error": {
    "code": 404,
    "message": "Invoice not found",
    "details": "Invoice with ID 999 does not exist"
  }
}
```

### Database Migrations

GORM auto-migrates schema on startup - no separate migration files needed for this simple app.

## API Endpoints

### Health Check

```shell
GET /api/health
```

Returns version info and service status.

### Invoices

```shell
GET    /api/invoices              # List all (paginated)
GET    /api/invoices/:id          # Get single invoice
POST   /api/invoices              # Create new invoice
PUT    /api/invoices/:id          # Update invoice
DELETE /api/invoices/:id          # Delete invoice
GET    /api/invoices/:id/pdf      # Download PDF
```

### Invoice Items (nested resource)

```shell
GET    /api/invoices/:id/items           # List items
POST   /api/invoices/:id/items           # Add item
PUT    /api/invoices/:id/items/:itemId   # Update item
DELETE /api/invoices/:id/items/:itemId   # Delete item
```

## Running Locally

### Development

```bash
cd backend

# Install dependencies
go mod download

# Run directly
go run main.go

# Or build binary
go build -o invoice-backend
./invoice-backend
```

Server starts on `http://localhost:8080`

### Environment Variables

```bash
source .env.example

# DB_PATH="/app/data/invoices.db" - Optional, defaults to this path
```

The app creates the database file automatically if it doesn't exist.

## Docker

### Build Image

```bash
docker build \
  --build-arg GO_VERSION=1.25 \
  --build-arg DEBIAN_VERSION=bookworm-20251020-slim \
  --build-arg APP_VERSION=1.0.0 \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  --build-arg BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -t invoice-backend .
```

Multi-stage build: Go for compilation, Debian slim for runtime.

### Run Container

```bash
docker run -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  invoice-backend
```

Mount `/app/data` to persist the SQLite database.

## Testing backend response

```bash
curl -X GET http://localhost:8080/api/health
```

## Security

- **Input validation** - Gin's binding tags validate JSON
- **SQL injection** - GORM uses prepared statements
- **Non-root user** - Dockerfile runs as UID 1000
- **No secrets** - no hardcoded credentials (use env vars)
