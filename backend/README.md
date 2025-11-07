# Backend Overview

Simple REST backend written in Go (Gin) that serves invoice data and PDF generation.

## Summary

- Small Gin-based HTTP API with MVC-like layout: routes → controllers → services → models.
- Files: main.go, routes/, controllers/, services/, models/, database/, config/.

## Technologies

- Go + Gin
- SQLite database (via database/connection.go)
- PDF generation service (services/pdf_service.go)
- Dockerfile for container builds

## Standards

- Semantic API responses (JSON)
- Linting expected in CI
- HTTP status codes for success/error handling (see below)
- Full REST

## What is required for full operation

- Go toolchain (recommended Go 1.25+)
- Environment variables (see `.env`)
- Docker for containerized runs

## Env setup (example)

- source .env
- build with go (`go build`)
- run binary

or

- `docker build -t invoice-backend .`
- `docker run --name invoice-backend -p 8080:8080 invoice-backend`

## REST overview (main endpoints)

- GET    /api/health               — health check
- GET    /api/invoices             — list invoices (with pagination/filter)
- POST   /api/invoices             — create invoice
- GET    /api/invoices/:id         — get invoice by id
- PUT    /api/invoices/:id         — update invoice
- DELETE /api/invoices/:id         — delete invoice
- POST   /api/invoices/:id/pdf     — generate/download PDF for invoice
- (Invoice items handled under /api/invoices/:id/items)

## Common HTTP status codes returned

- 200 OK — successful GET/PUT
- 201 Created — resource created
- 204 No Content — successful delete (no body)
- 400 Bad Request — validation or client error
- 404 Not Found — missing resource
- 500 Internal Server Error — unexpected server/database error

## Structure notes

- routes/routes.go wires endpoints to controllers.
- controllers handle request/response and call models/services.
- database/connection.go manages DB connection pooling.
- services/pdf_service.go handles PDF creation.
