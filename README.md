
# Simple GO invoice

A simple single-page application (SPA) demonstrating a basic full-stack setup using **Go (Gin)** for the backend and **Vue 3 + Vite + Vuetify** for the frontend.  

The project serves as a lightweight example of a RESTful API communicating with a modern JavaScript frontend.

---

## Deployment Status

[![CI Pipeline](https://github.com/JustFiesta/simple-go-invoice/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/JustFiesta/simple-go-invoice/actions/workflows/ci.yaml)
[![Terraform Apply](https://github.com/JustFiesta/simple-go-invoice/actions/workflows/terraform-apply.yaml/badge.svg)](https://github.com/JustFiesta/simple-go-invoice/actions/workflows/terraform-apply.yaml)
[![Security Scanning](https://github.com/JustFiesta/simple-go-invoice/actions/workflows/security-scans.yaml/badge.svg)](https://github.com/JustFiesta/simple-go-invoice/actions/workflows/security-scans.yaml)
---

## Stack

- **Backend:** Go + Gin (REST API)
- **Frontend:** Vue 3 + Vite + Vuetify
- **Containerization:** Docker + Docker Compose
- **IaC**: Terraform, K8 + Kustomize
- **CI/CD**: Github Acitons with composite Actions, ArgoCD on K8, Annual Security Scanning with Trivy

---

## Features

- Simple REST API
- Vue SPA with Vuetify UI components
- Sort/Search/Filter/Pagination
- Communication between frontend and backend via HTTP
- Docker + Compose setup for local development
- Treafig reverse proxy for one URL to backend/frontend
- Kubernetes Deployment
- Full CI/CD

---

## Running the Project

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/invoice-spa-go.git

cd invoice-spa-go
```

### 2. Copy .env.example

```shell
cp .env.example .env
```

### 2. Start/Stop with Docker Compose

```shell
docker compose up --build -d .

docker compose down -v # -v is optional as it will remove volumes
```

#### Manual build

```shell
docker build --build-arg GO_VERSION=$GO_VERSION -t invoice-backend ./backend
docker build --build-arg NODE_VERSION=$NODE_VERSION -t invoice-frontend ./frontend
```

#### Debug

```shell
docker compose watch
```

### 3. Browser access

```txt
Frontend: http://localhost:8000
Backend API: http://localhost:8000/api/hello
```

---

## Notes

The project is meant for educational or demo purposes.

It is ment to showcase usage of DevOps principals in simple manner.

---

## Tool docs reference

| Tool          | Description                        | Link                                                                                   |
| ------------------ | --------------------------- | -------------------------------------------------------------------------------------- |
| **Gin (Go)**       | HTTP microframework for Go  | [https://gin-gonic.com/docs/](https://gin-gonic.com/docs/)                             |
| **Vue 3**          | Frontend framework JS       | [https://vuejs.org/guide/introduction.html](https://vuejs.org/guide/introduction.html) |
| **Vite**           | Bundler and dev server | [https://vitejs.dev/guide/](https://vitejs.dev/guide/)                                 |
| **Vuetify**        | Ready UI components        | [https://next.vuetifyjs.com/](https://next.vuetifyjs.com/)                             |
| **Docker Compose** | Containers definition        | [https://docs.docker.com/compose/](https://docs.docker.com/compose/)                   |
