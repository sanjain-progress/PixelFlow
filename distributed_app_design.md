# PixelFlow: Distributed Image Processing Pipeline

## 1. Project Idea
**PixelFlow** is a distributed system where users upload images to be processed (resized, filtered, watermarked) asynchronously. It mimics real-world heavy processing jobs.

**Why this fits:**
- **Microservices**: Separate **Auth**, **API (Gateway/Task)**, and **Worker** services.
- **Inter-Service Comm**: **gRPC** for synchronous calls (API -> Auth), **Kafka** for async events (API -> Worker).
- **Polyglot Persistence**: **PostgreSQL** for Users (Relational), **MongoDB** for Tasks (Document).
- **Scalability**: Independent scaling of all 3 components.

---

## 2. High-Level Architecture

```ascii
                                      +----------------+
                                      |  Grafana/Loki  |
                                      +-------+--------+
                                              ^
+-------------+                               | (Metrics/Logs)
| React UI    |                               |
| (Frontend)  |                               |
+------+------+                               |
       | (HTTP/JSON)                          |
       v                                      |
+-------------+       (gRPC)          +-------+--------+
| NGINX       | --------------------> | Auth Service   |
| (Ingress)   | <-------------------- | (Go + gRPC)    |
+------+------+                       +-------+--------+
       |                                      |
       | (HTTP /api/tasks)                    v
       v                              +----------------+
+-------------+                       | PostgreSQL     |
| API Service |                       | (Users DB)     |
| (Go + HTTP) |                       +----------------+
+------+------+
       |
       | (Produce: "ImageUploaded")
       v
+-------------+    +-------------+
| Kafka       |    | MongoDB     |
| (Broker)    |    | (Tasks DB)  |
+------+------+    +-------------+
       |
       | (Consume: "ImageUploaded")
       v
+-------------+
| Go Worker   | <---- (Scale to N replicas)
| (Consumer)  |
+------+------+
       |
       v
+-------------+
| Object Store|
| (MinIO/S3)  |
+-------------+
```

---

## 3. Tech-by-Tech Usage Plan

| Technology | Usage in PixelFlow |
| :--- | :--- |
| **Go (Golang)** | **Auth Service**: gRPC server, handles Login/Signup, issues JWTs. <br> **API Service**: HTTP REST API, validates JWTs (via gRPC check or local key), manages Tasks. <br> **Worker Service**: Consumes Kafka, processes images. |
| **gRPC / Protobuf** | **Communication**: API Service calls Auth Service to validate tokens or fetch user details. |
| **PostgreSQL** | **Auth DB**: Stores User accounts (Relational data). Best practice for structured user data. |
| **MongoDB** | **Task DB**: Stores Task metadata (Flexible schema). |
| **Kafka** | **Async Messaging**: Decouples API from Workers. |
| **Docker & K8s** | Containerization and Orchestration for all 3 services + databases. |
| **Helm** | Deploys the entire stack (Auth, API, Worker, DBs). |
| **NGINX Ingress** | Routes `/auth` to Auth Service (optional, or API Gateway pattern) and `/api` to API Service. |
| **Prometheus/Grafana** | Metrics for "gRPC Latency", "HTTP Latency", "Kafka Lag". |
| **Loki** | Distributed logging across all 3 services. |

---

## 4. Feature Set

1.  **Auth Service**:
    -   `Register(username, password)` -> Saves to Postgres.
    -   `Login(username, password)` -> Returns JWT.
    -   `Validate(token)` -> gRPC method called by API Service.
2.  **API Service**:
    -   **Middleware**: Intercepts requests, calls Auth Service (via gRPC) to validate identity.
    -   `POST /tasks`: Upload image -> Save to MinIO -> Save to Mongo -> Push to Kafka.
    -   `GET /tasks`: Fetch user's tasks from Mongo.
3.  **Worker Service**:
    -   Consumes `task_created`.
    -   Processes image.
    -   Updates Mongo `status="completed"`.

---

## 5. Monitoring Metrics

**New Microservice Metrics:**
-   `grpc_server_handled_total`: Auth Service request count.
-   `grpc_server_handling_seconds`: Auth Service latency.

---

## 6. Step-by-Step Learning Roadmap

### Phase 1: The Foundation (Auth Service)
-   **Goal**: Build the Identity Provider.
-   **Tech**: Go, gRPC, Protobuf, PostgreSQL.
-   **Tasks**:
    -   Define `auth.proto`.
    -   Implement Go gRPC Server.
    -   Connect to PostgreSQL.
    -   Implement JWT generation.

### Phase 2: The Gateway (API Service)
-   **Goal**: Build the main REST API.
-   **Tech**: Go, HTTP, MongoDB, gRPC Client.
-   **Tasks**:
    -   Create REST endpoints.
    -   Implement gRPC Client to talk to Auth Service.
    -   Connect to MongoDB.

### Phase 3: The Async Engine (Kafka + Worker)
-   **Goal**: Background processing.
-   **Tech**: Kafka, Go.
-   **Tasks**:
    -   Setup Kafka.
    -   API produces events.
    -   Worker consumes and processes images.

### Phase 4-7: (Same as before - UI, Docker, K8s, CI/CD)
-   React UI will now login against the API (which proxies to Auth) or directly to Auth if exposed.

---
