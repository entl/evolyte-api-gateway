# API Gateway

![Go](https://img.shields.io/badge/Go-1.23-blue?logo=go&logoColor=white)
![Echo](https://img.shields.io/badge/Echo_Framework-Web-blue?logo=go)
![SQLC](https://img.shields.io/badge/sqlc-SQL%20Codegen-4B8BBE?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-Cache-DC382D?logo=redis&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-Monitoring-orange?logo=prometheus)
![Grafana](https://img.shields.io/badge/Grafana-Dashboard-F46800?logo=grafana)
![Docker](https://img.shields.io/badge/Docker-Container-2496ED?logo=docker)
![Elasticsearch](https://img.shields.io/badge/Elasticsearch-Search-005571?logo=elasticsearch)
![Kibana](https://img.shields.io/badge/Kibana-Logs-005571?logo=kibana)

---
A lightweight reverse-proxy written in Go that provides centralized **JWT authentication**, **role-based header injection**, and **response caching** for GET requests via Redis.  
It sits in front of micro-services, validating each request and forwarding it to the appropriate backend while shielding internal endpoints from direct exposure.

---

## ✨ Key Features

| Feature | Description |
|---------|-------------|
| **JWT Authentication** | Verifies the `Authorization: Bearer <token>` header for validity. |
| **Role Extraction** | Parses the token’s `roles` (or configurable claim) and injects them into the upstream request as `X-User-Roles`. |
| **Request Proxying** | Transparently forwards the original method, path, query string, and body to the target service defined in `config.yaml`. |
| **Redis-backed Response Cache** | Automatically caches successful **GET** responses with default **TTL** of 60 seconds. Subsequent identical requests are served from Redis, cutting latency and offloading the service. |
| **Config-Driven Routing** | Adds, removes, or secures routes without recompiling—just edit the YAML file and restart. |
| **Prometheus Metrics** | Exposes real-time metrics at `/metrics` endpoint, compatible with Prometheus scraping. |
| **Grafana Dashboards** | Includes a ready-to-use `grafana-dashboard.json` for visualising request rates, latencies, and errors in Grafana. |

---

## 🏗️ Architecture

```
            ┌──────────┐
            │  Grafana │
            └──────────┘
                  ▲
                  │
           ┌──────────────┐      ┌────────┐
 Client ───►    Gateway   ├─────►│Service │
           └──────────────┘      └────────┘
                 ▲   │
┌───────────┐    │   │
│Redis Cache│  ◄─┘   └─► Roles headers (X-User-Roles, X-User-ID)
└───────────┘ 
```

---

## 📋 Configuration

```yaml
# config.yaml

services:
  main:
    backend: "http://api:8001"
    routes:
      - from_path: "/api/v1/{prefix}/health"
        to_path: "/api/v1/health"
        method: "GET"
        auth_required: true
        allowed_roles: ["admin", "user"]
      ...
  other_service:
    backend: "http://other-service:8000"
    routes:
      - from_path: "/api/v1/{prefix}/model-info"
        to_path: "/api/v1/model-info"
        method: "GET"
        auth_required: false
```

---

## 🔧 Route Configuration Explained

| Field          | Description |
|----------------|-------------|
| `from_path`     | Public API path received by the gateway. Wildcards (`*`) supported. |
| `to_path`       | Internal path forwarded to backend service. Wildcards must match `from_path`. |
| `method`        | HTTP method to match (`GET`, `POST`, `PATCH`, `DELETE`). |
| `auth_required` | If true, JWT is required and validated. |
| `allowed_roles` | Optional list of roles permitted to access this route. |

---


## 🚀 Running Locally

```bash
docker compose -f docker-compose.dev.yml --env-file .env up --build
```

---

## 📈 Metrics & Observability

The gateway exposes metrics via `echo-prometheus-middleware` at `/metrics`.

### Prometheus

Prometheus scrapes the metrics via Docker Compose and is available at `http://localhost:9090`. Example config:

```yaml
scrape_configs:
  - job_name: 'gateway'
    static_configs:
      - targets: ['gateway:8080']
```

### Grafana

Grafana visualizes these metrics at `http://localhost:3000`.

To import the dashboard:

1. Open Grafana
2. Go to **Dashboards → Import**
3. Upload `grafana-dashboard.json`
4. Select Prometheus as the data source

---

## 📑 License

MIT © 2025 Maksym Vorobyov
