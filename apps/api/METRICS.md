# API Service Metrics

The API Service exposes Prometheus metrics at `/metrics`.

## HTTP Metrics

| Metric Name | Type | Description | Labels |
|---|---|---|---|
| `api_requests_total` | Counter | Total number of HTTP requests | `method`, `endpoint`, `status` |
| `api_request_duration_seconds` | Histogram | Request latency in seconds | `method`, `endpoint` |

## Business Metrics

| Metric Name | Type | Description |
|---|---|---|
| `api_tasks_created_total` | Counter | Total number of tasks created via `/api/upload` |
| `api_tasks_retrieved_total` | Counter | Total number of task list requests via `/api/tasks` |

## Kafka Metrics

| Metric Name | Type | Description |
|---|---|---|
| `api_kafka_messages_published_total` | Counter | Total messages successfully published to Kafka |
| `api_kafka_publish_errors_total` | Counter | Total errors when publishing to Kafka |

## Example Queries

### Request Rate
```promql
rate(api_requests_total[5m])
```

### Task Creation Rate
```promql
rate(api_tasks_created_total[5m])
```

### Kafka Publish Success Rate
```promql
rate(api_kafka_messages_published_total[5m])
```

### Kafka Publish Error Rate
```promql
rate(api_kafka_publish_errors_total[5m])
```
