# API Service Implementation Plan

## Goal
Handle image processing tasks, file uploads, and coordinate with the Worker service via Kafka.

## Architecture
- **Type**: REST API
- **Port**: 8080
- **Database**: MongoDB (Task metadata)
- **Message Queue**: Kafka (Producer)
- **Auth**: Client of Auth Service (via HTTP)

## API Endpoints

### POST /api/upload
- **Auth**: Required
- **Input**: `image_url`
- **Logic**:
  1. Validate token via Auth Service
  2. Create Task record in MongoDB (Status: PENDING)
  3. Publish "TaskCreated" event to Kafka
  4. Return Task ID

### GET /api/tasks
- **Auth**: Required
- **Logic**:
  1. Validate token via Auth Service
  2. Query MongoDB for tasks belonging to user
  3. Return list of tasks

## Internal Components
- **middleware**: Auth middleware (calls Auth Service)
- **kafka**: Producer for publishing events
- **db**: MongoDB connection
- **models**: Task struct definition

## Dependencies
- `github.com/gin-gonic/gin`
- `go.mongodb.org/mongo-driver`
- `github.com/confluentinc/confluent-kafka-go/kafka`
