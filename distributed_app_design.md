# PixelFlow - Distributed Image Processing Application

## Overview
PixelFlow is a production-style distributed application for asynchronous image processing, built to demonstrate modern backend, cloud-native, and microservices patterns.

**Status**: ✅ **Phase 3 Complete - All Services Deployed and Tested**

## Architecture

### High-Level Design (HTTP + Event-Driven)

```
┌──────────┐       HTTP       ┌──────────┐       HTTP        ┌──────────┐       HTTP        ┌──────────┐
│   User   │ ───────────────▶ │ Frontend │ ────────────────▶ │   API    │ ────────────────▶ │   Auth   │
│ (Browser)│                  │ Service  │                   │ Service  │                   │ Service  │
└──────────┘                  └──────────┘                   └────┬─────┘                   └──────────┘
                                                                  │                               │
                                                                  │                               ▼
                                                                  │                         ┌──────────┐
                                                                  │                         │PostgreSQL│
                                                                  │                         │  (Users) │
                                                                  ▼                         └──────────┘
                                                            ┌──────────┐
                                                            │  MongoDB │
                                                            │  (Tasks) │
                                                            └──────────┘
                                                                  │
                                                                  ▼
                                                            ┌──────────┐
                                                            │  Kafka   │
                                                            │ (Events) │
                                                            └────┬─────┘
                                                                 │
                                                                 │ Subscribe
                                                                 ▼
                                                           ┌──────────┐
                                                           │  Worker  │
                                                           │ Service  │
                                                           └────┬─────┘
                                                                 │
                                                                 ▼
                                                           ┌──────────┐
                                                           │  MongoDB │
                                                           │ (Update) │
                                                           └──────────┘
```

## Services

### 1. Frontend Service (UI)
**Purpose**: User Interface for the application
**Port**: 3000
**Type**: Single Page Application (SPA)
**Tech Stack**: React, TailwindCSS, Nginx (Docker)

**Features**:
- User Registration & Login pages
- Dashboard for task management
- Image URL upload form
- Real-time task status updates (Polling)
- Protected route management (JWT)

### 2. Auth Service (HTTP)
**Purpose**: User authentication and JWT token management  
**Port**: 50051  
**Protocol**: HTTP REST  
**Database**: PostgreSQL

**Endpoints**:
- `POST /register` - Create new user account
- `POST /login` - Authenticate and receive JWT token
- `GET /validate` - Validate JWT token (used by API service)

**Tech Stack**:
- Go with Gin framework
- PostgreSQL with GORM ORM
- JWT for token generation
- bcrypt for password hashing

### 3. API Service (HTTP)
**Purpose**: REST API for task management  
**Port**: 8080  
**Protocol**: HTTP REST  
**Databases**: MongoDB (tasks), Kafka (events)

**Endpoints**:
- `GET /health` - Health check (public)
- `POST /api/upload` - Create image processing task (authenticated)
- `GET /api/tasks` - List user's tasks (authenticated)

**Authentication**: JWT validation via HTTP call to Auth Service

**Tech Stack**:
- Go with Gin framework
- MongoDB for task storage
- Kafka producer for event publishing
- HTTP client for auth validation

### 4. Worker Service (Background)
**Purpose**: Asynchronous image processing  
**Protocol**: Kafka consumer  
**Databases**: MongoDB (task updates)

**Functionality**:
- Consumes tasks from Kafka topic `image-tasks`
- Simulates image processing (5-second delay)
- Updates task status: PENDING → PROCESSING → COMPLETED
- Generates processed image URL

**Tech Stack**:
- Go
- Kafka consumer (consumer group: `worker-group-1`)
- MongoDB for status updates

## Data Flow

### 1. User Registration & Login
```
Client → POST /register → Auth Service → PostgreSQL (create user)
Client → POST /login → Auth Service → PostgreSQL (verify) → JWT Token
```

### 2. Task Creation
```
Client → POST /api/upload (with JWT) → API Service
  ↓
API validates JWT via Auth Service HTTP call
  ↓
API creates task in MongoDB (status: PENDING)
  ↓
API publishes event to Kafka topic
  ↓
Returns task to client
```

### 3. Background Processing
```
Kafka → Worker Service (consumer)
  ↓
Worker updates MongoDB (status: PROCESSING)
  ↓
Worker simulates processing (5 seconds)
  ↓
Worker updates MongoDB (status: COMPLETED, processed_url)
```

### 4. Status Check
```
Client → GET /api/tasks (with JWT) → API Service
  ↓
API validates JWT
  ↓
API queries MongoDB
  ↓
Returns tasks with current status
```

## Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Languages** | Go 1.21+ | All services |
| **API Framework** | Gin | HTTP REST APIs |
| **Auth Database** | PostgreSQL 15 | User accounts |
| **Task Database** | MongoDB 6.0 | Task storage |
| **Message Queue** | Kafka 7.3 | Event streaming |
| **ORM** | GORM | PostgreSQL access |
| **Containerization** | Docker, Docker Compose | Deployment |

## Infrastructure (Docker)

All services run in Docker containers:
- `pixelflow-frontend` - React Frontend (Nginx)
- `pixelflow-auth` - Auth Service
- `pixelflow-api` - API Service  
- `pixelflow-worker` - Worker Service
- `pixelflow-postgres-auth` - PostgreSQL database
- `pixelflow-mongo` - MongoDB database
- `pixelflow-kafka` - Kafka broker
- `pixelflow-zookeeper` - Kafka coordination

## Development vs Original Design

**Original Plan**: gRPC for inter-service communication  
**Implemented**: HTTP REST for simplicity and reduced complexity

**Benefits of HTTP Approach**:
- ✅ Simpler build process (no protobuf code generation)
- ✅ Easier debugging with standard HTTP tools
- ✅ Lower learning curve
- ✅ Still demonstrates microservices patterns
- ✅ Maintains distributed architecture principles

## Running the Application

```bash
# Start all services
make up

# Check status
make ps

# Access UI
# Open http://localhost:3000

# Run E2E tests
make test

# View logs
make logs            # All services
make logs-worker     # Specific service

# Stop all services
make down
```

For all available commands, run `make help`.

## Verified Features ✅

- ✅ **Frontend UI**: Login, Register, Dashboard, Task Upload
- ✅ User registration and authentication
- ✅ JWT token generation and validation  
- ✅ Protected API endpoints
- ✅ Task creation and persistence
- ✅ Kafka event publishing
- ✅ Worker consumption and processing
- ✅ Status updates (PENDING → PROCESSING → COMPLETED)
- ✅ End-to-end workflow tested

## Key Learning Outcomes

### Backend Development
- RESTful API design with Gin
- JWT authentication patterns
- Database modeling (PostgreSQL + MongoDB)
- ORM usage with GORM

### Distributed Systems
- Microservices architecture
- Event-driven design with Kafka
- Service-to-service communication (HTTP)
- Async task processing patterns

### DevOps & Cloud-Native
- Docker containerization
- Multi-container orchestration (Docker Compose)
- Service discovery and networking
- Health checks and logging

### Scalability Patterns
- Horizontal scaling with Kafka consumer groups
- Database per service pattern
- Stateless service design
- Message queue for decoupling

---
