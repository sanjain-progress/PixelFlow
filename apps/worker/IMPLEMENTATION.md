# Worker Service Implementation Plan

## Goal
Process background image tasks asynchronously by consuming Kafka messages.

## Architecture
- **Type**: Background Worker (Daemon)
- **Database**: MongoDB (Update task status)
- **Message Queue**: Kafka (Consumer)

## Logic Flow

1. **Consume**: Listen to `image-tasks` topic
2. **Process**:
   - Parse message (Task ID, Image URL)
   - Simulate image processing (sleep)
   - Generate "processed" URL
3. **Update**:
   - Update Task status in MongoDB to "PROCESSING"
   - Complete processing
   - Update Task status in MongoDB to "COMPLETED" with `processed_url`

## Internal Components
- **kafka**: Consumer group implementation
- **processor**: Image processing logic
- **db**: MongoDB connection
- **models**: Task struct definition

## Dependencies
- `go.mongodb.org/mongo-driver`
- `github.com/confluentinc/confluent-kafka-go/kafka`
