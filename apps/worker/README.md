# Worker Service
 
 The Worker Service is a background consumer that processes image tasks. It listens to Kafka topics, simulates processing, and updates task status in MongoDB.
 
 ## 🏗️ Architecture
 
 ```mermaid
 stateDiagram-v2
     [*] --> Idle
     Idle --> Consuming: New Message
     Consuming --> Processing: Parse Event
     Processing --> UpdatingDB: Process Image (5s)
     UpdatingDB --> Idle: Update Status (COMPLETED)
 ```
 
 ## 🔄 Workflow
 1. Consume message from `image-tasks` topic.
 2. Parse JSON payload (Task ID, Image URL).
 3. Simulate processing (sleep 5s).
 4. Update MongoDB document status to `COMPLETED`.
 
 ## 🛠️ Tech Stack
 - **Language**: Go
 - **Messaging**: Kafka (Consumer Group)
 - **Database**: MongoDB
 - **Observability**: Prometheus Metrics, Jaeger Tracing
