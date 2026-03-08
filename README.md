# URL Shortener

A multi-tenant URL shortening platform with real-time analytics, built as a microservices monorepo.

## Architecture

```
                                 ┌──────────────────────┐
                                 │   Angular Frontend   │
                                 │   (url-shortener-web) │
                                 └───────┬──────┬───────┘
                                         │      │
                              ┌──────────▼─┐  ┌─▼──────────┐
                              │ URL Short. │  │  Dashboard  │
                              │    API     │  │   Service   │
                              │ Go :8080   │  │ Java :8083  │
                              └──────┬─────┘  └──────┬──────┘
                                     │               │
      ┌──────────────┐         ┌─────▼─────┐   ┌────▼──────┐
      │  Auth Service │         │  MongoDB  │   │ClickHouse │
      │   Go :8081    │         └─────┬─────┘   └───────────┘
      └──────────────┘               │                ▲
                              ┌──────▼──────┐         │
┌────────────┐  Kafka   ┌─────┤ Redirector  │   ┌─────┴──────┐
│  Browser   ├─────────►│     │  Go :8082   ├──►│  Analytics │
│ GET /:code │◄─── 307 ─┤     └─────────────┘   │   (Flink)  │
└────────────┘          │                        └────────────┘
                        └──► Redis
```

## Services

| Service | Language | Port | Description |
|---|---|---|---|
| **url-shortener-api** | Go | 8080 | Creates short URLs, manages URL mappings |
| **auth-service** | Go | 8081 | JWT-based authentication |
| **redirector-service** | Go | 8082 | Resolves short codes → 307 redirect, publishes click events to Kafka |
| **dashboard-service** | Java/Spring Boot | 8083 | REST API for analytics data (queries ClickHouse) |
| **analytics-service** | Java/Flink | — | Stream processing pipeline: Kafka → ClickHouse |
| **url-shortener-web** | Angular | 4200 | Frontend SPA |

## Infrastructure

| Component | Purpose |
|---|---|
| **MongoDB** | URL mappings, user data |
| **Redis** | URL cache (24h TTL), global counter for short code generation |
| **Kafka** | Click event streaming (`click-events` topic) |
| **ClickHouse** | Analytics OLAP — raw events + materialized views for hourly/daily/weekly aggregation |
| **Zookeeper** | Kafka coordination |

## Project Layout

```
url-shortner/
├── apps/
│   ├── url-shortener-api/        # Go — URL shortening API
│   │   ├── controller/           #   HTTP handlers
│   │   ├── service/              #   Business logic
│   │   ├── repository/           #   MongoDB + Redis cache
│   │   ├── model/                #   Data models
│   │   ├── middleware/           #   JWT authentication
│   │   ├── infrastructure/       #   DB/Redis/Auth clients
│   │   ├── utils/counter/        #   Redis-based global counter
│   │   ├── utils/hashing/        #   Base62 encoding
│   │   └── bootstrap/            #   Wire DI + router setup
│   │
│   ├── redirector-service/       # Go — Short URL redirect
│   │   ├── controller/           #   Redirect handler
│   │   ├── service/              #   Redirect + click event publishing
│   │   ├── model/                #   ClickEvent model
│   │   └── infrastructure/       #   Kafka producer
│   │
│   ├── analytics-service/        # Java/Flink — Stream processing
│   │   └── src/main/java/
│   │       ├── sink/             #   ClickHouse JDBC sink
│   │       ├── deserializer/     #   Kafka JSON deserializer
│   │       └── model/            #   ClickEvent POJO
│   │
│   ├── dashboard-service/        # Java/Spring Boot — Analytics API
│   │   └── src/main/java/
│   │       ├── controller/       #   REST endpoints
│   │       ├── service/          #   ClickHouse + MongoDB queries
│   │       ├── repository/       #   Reactive MongoDB repos
│   │       ├── model/            #   MongoDB documents
│   │       └── dto/              #   Response DTOs
│   │
│   ├── auth-service/             # Go — Authentication
│   └── url-shortener-web/        # Angular — Frontend
│
└── ops/config/
    ├── compose.yaml              # Docker Compose for all infra
    └── clickhouse/init.sql       # ClickHouse schema + materialized views
```

## Getting Started

### Prerequisites

- **Go 1.21+**
- **Java 21** (JDK)
- **Node.js 18+** (for Angular frontend)
- **Docker & Docker Compose**

### 1. Start Infrastructure

```bash
docker compose -f ops/config/compose.yaml up -d
```

This starts MongoDB, Redis, Kafka, Zookeeper, and ClickHouse. The ClickHouse schema (tables + materialized views) is auto-created on first start.

### 2. Start Backend Services

```bash
# URL Shortener API (port 8080)
cd apps/url-shortener-api && go run main.go

# Auth Service (port 8081)
cd apps/auth-service && go run main.go

# Redirector Service (port 8082)
cd apps/redirector-service && go run main.go

# Analytics Flink Pipeline
cd apps/analytics-service && ./gradlew shadowJar
java -jar build/libs/analytics-service-0.0.1-SNAPSHOT-all.jar

# Dashboard Service (port 8083)
cd apps/dashboard-service && ./gradlew bootRun
```

### 3. Start Frontend

```bash
cd apps/url-shortener-web && npm install && ng serve
```

## API Reference

### URL Shortener API (`:8080`)

All endpoints require a JWT token in the `Authorization` header.

#### Create Short URLs

```
POST /api/v1/urls/shorten
```

**Request:**
```json
{
  "links": [
    { "url": "https://example.com/very/long/path" },
    { "url": "https://another-site.com/page" }
  ]
}
```

**Response:**
```json
[
  {
    "long_url": "https://example.com/very/long/path",
    "short_url": "http://localhost:8080/a1b2c3"
  },
  {
    "long_url": "https://another-site.com/page",
    "short_url": "http://localhost:8080/d4e5f6"
  }
]
```

#### List URL Mappings

```
GET /api/v1/urls
```

Returns all URLs created by the authenticated user's tenant.

### Redirector Service (`:8082`)

```
GET /:code → 307 Redirect to original URL
```

Each redirect publishes a click event to Kafka for analytics.

### Dashboard Service — Analytics API (`:8083`)

| Endpoint | Description |
|---|---|
| `GET /api/v1/analytics/{tenantId}/summary` | Total clicks, unique IPs, unique short codes |
| `GET /api/v1/analytics/{tenantId}/clicks?page=0&size=20` | Paginated recent click events |
| `GET /api/v1/analytics/{tenantId}/top-links?limit=10&granularity=HOUR` | Most clicked links |
| `GET /api/v1/analytics/{tenantId}/timeseries?granularity=DAY` | Clicks over time |
| `GET /api/v1/analytics/{tenantId}/referers` | Traffic source breakdown |

**Granularity options:** `HOUR`, `DAY`, `WEEK`

## Analytics Pipeline

```
Redirect Request
       │
       ▼
  Kafka Topic: click-events
       │
       ▼
  Flink Pipeline (analytics-service)
       │ JDBC insert
       ▼
  ClickHouse: url_shortener.click_events
       │ auto-populated by materialized views
       ├──▶ click_aggregates_hourly
       ├──▶ click_aggregates_daily
       └──▶ click_aggregates_weekly
```

**Click Event Schema:**

| Field | Type | Description |
|---|---|---|
| `short_code` | String | The short URL code |
| `original_url` | String | Target URL |
| `tenant_id` | String | Tenant identifier |
| `timestamp` | DateTime64 | Event time (UTC) |
| `user_agent` | String | Browser user agent |
| `referer` | String | HTTP referer header |
| `ip_address` | String | Client IP address |

## Tech Stack

| Layer | Technology |
|---|---|
| API Gateway / URL Shortener | **Go** + Gin |
| Redirect Service | **Go** + Gin + kafka-go |
| Authentication | **Go** + JWT (RS256) |
| Stream Processing | **Apache Flink** 1.20 (Java) |
| Analytics API | **Spring Boot** 4.0 + WebFlux |
| Frontend | **Angular** |
| Primary Database | **MongoDB** 8 |
| Cache & Counters | **Redis** 7 |
| Event Streaming | **Apache Kafka** (Confluent 7.6) |
| Analytics OLAP | **ClickHouse** 24.12 |
| DI (Go) | **Google Wire** |
| DI (Java) | **Spring IoC** |
