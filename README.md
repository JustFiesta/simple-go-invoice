# Simple GO invoice

A simple single-page application (SPA) demonstrating a basic full-stack setup using **Go (Gin)** for the backend and **Vue 3 + Vite + Vuetify** for the frontend.  
The project serves as a lightweight example of a RESTful API communicating with a modern JavaScript frontend.

---

## Stack

- **Backend:** Go + Gin (REST API)
- **Frontend:** Vue 3 + Vite + Vuetify
- **Containerization:** Docker + Docker Compose

---

## Features

- Simple REST API
- Vue SPA with Vuetify UI components
- Sort/Search/Filter/Pagination
- Communication between frontend and backend via HTTP
- Docker + Compose setup for local development
- Treafig reverse proxy for one URL to backend/frontend

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

Suitable as a base for simple CRUD applications (e.g., invoice or task manager systems).

---

## Tool docs reference

| Tool          | Description                        | Link                                                                                   |
| ------------------ | --------------------------- | -------------------------------------------------------------------------------------- |
| **Gin (Go)**       | Mikroframework HTTP dla Go  | [https://gin-gonic.com/docs/](https://gin-gonic.com/docs/)                             |
| **Vue 3**          | Frontend framework JS       | [https://vuejs.org/guide/introduction.html](https://vuejs.org/guide/introduction.html) |
| **Vite**           | Szybki bundler i dev server | [https://vitejs.dev/guide/](https://vitejs.dev/guide/)                                 |
| **Vuetify**        | Gotowe komponenty UI        | [https://next.vuetifyjs.com/](https://next.vuetifyjs.com/)                             |
| **Docker Compose** | Definicja kontenerów        | [https://docs.docker.com/compose/](https://docs.docker.com/compose/)                   |
