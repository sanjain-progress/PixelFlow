# API Service
 
 The API Service manages image processing tasks. It handles file uploads (simulated), creates task records in MongoDB, and publishes events to Kafka.
 
 ## 🏗️ Architecture
 
 ```mermaid
 sequenceDiagram
     participant User
     participant API
     participant Auth
     participant DB
     participant Kafka
 
     User->>API: POST /api/upload (Header: Bearer Token)
     API->>Auth: GET /validate
     Auth-->>API: Token Valid (User ID)
     API->>DB: Create Task (Status: PENDING)
     API->>Kafka: Publish Event (image-tasks)
     API-->>User: Task Created (201 Created)
 ```
 
 ## 🚀 API Endpoints
 
 | Method | Endpoint | Description |
 |--------|----------|-------------|
 | GET | `/health` | Service health check |
 | POST | `/api/upload` | Create a new task |
 | GET | `/api/tasks` | List all tasks for user |
 | GET | `/metrics` | Prometheus metrics |
 
 ## 🛠️ Tech Stack
 - **Framework**: Gin
 - **Database**: MongoDB
 - **Messaging**: Kafka (Producer)
 - **Auth**: JWT Middleware
