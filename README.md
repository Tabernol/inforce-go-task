# Rarible Service (Go)

A Go service and client to interact with the Rarible API.
The service exposes HTTP endpoints to query NFT ownerships and trait rarities.

---

## Features

* **Rarible Client**

  * Get NFT ownership by ID
  * Query NFT trait rarities
* **HTTP Service**

  * `/health` — health check
  * `/api/v1/rarible/ownership?ownershipId=...` — get NFT ownership
  * `/api/v1/rarible/traits` — query trait rarities (POST)
* Fully tested (unit and handler tests)
* Runs in Docker

---

## Prerequisites

* Go >= 1.23
* Docker & Docker Compose

---

## Getting Started

### 1. Clone repository

```bash
git clone git@github.com:Tabernol/inforce-go-task.git
cd inforce-go-task
```

### 2. Environment Variables

Create a `.env` file in the root:

```dotenv
RARIBLE_BASE_URL=https://api.rarible.org
RARIBLE_API_KEY=your_api_key_here
RARIBLE_TIMEOUT=10s
```

---

### 3. Run Locally

```bash
go run ./cmd
```

The service will start on `http://localhost:8080`.

---

### 4. Run in Docker

Build the Docker image:

```bash
docker build -t rarible-service .
```

Run container:

```bash
docker run --rm -p 8080:8080 --env-file .env rarible-service
```

Healthcheck:

```bash
curl http://localhost:8080/health
```

---

### 5. Run with Docker Compose

```bash
docker-compose up --build
```

---

### 6. API Usage

#### Get Ownership

```http
GET /api/v1/rarible/ownership?ownershipId=ETHEREUM:<contract>:<tokenId>:<owner>
```

#### Query Trait Rarity

```http
POST /api/v1/rarible/traits
Content-Type: application/json

{
  "collectionId": "ETHEREUM:0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d",
  "properties": [
    {"key": "background", "value": "blue"}
  ],
  "limit": 10
}
```

---

### 7. Testing

Run unit and handler tests:

```bash
go test ./... -v
```

> Tests cover the Rarible client and HTTP handlers.

---

## Project Structure

```
├── cmd/                   # Main entrypoint
├── internal/config        # Configuration loader
├── internal/server        # HTTP handlers and routing
├── pkg/rarible            # Rarible client and models
├── Dockerfile
└── docker-compose.yml
```

---
