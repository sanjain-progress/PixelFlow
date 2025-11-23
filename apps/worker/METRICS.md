# Worker Service Metrics

The Worker Service exposes Prometheus metrics at `/metrics` on port `8081`.

## Task Processing Metrics

| Metric Name | Type | Description | Labels |
|---|---|---|---|
| `worker_tasks_processed_total` | Counter | Total number of tasks processed | `status` (success/failure) |
| `worker_task_processing_duration_seconds` | Histogram | Duration of task processing in seconds | - |
| `worker_active_processing_tasks` | Gauge | Number of tasks currently being processed | - |

## Kafka Metrics

| Metric Name | Type | Description |
|---|---|---|
| `worker_kafka_messages_consumed_total` | Counter | Total messages consumed from Kafka |
| `worker_kafka_consumption_errors_total` | Counter | Total errors when consuming from Kafka |

## Example Queries

### Processing Rate
```promql
rate(worker_tasks_processed_total[5m])
```

### Success Rate
```promql
sum(rate(worker_tasks_processed_total{status="success"}[5m])) 
/ 
sum(rate(worker_tasks_processed_total[5m])) * 100
```

### Average Processing Duration
```promql
rate(worker_task_processing_duration_seconds_sum[5m]) 
/ 
rate(worker_task_processing_duration_seconds_count[5m])
```

### Active Tasks
```promql
worker_active_processing_tasks
```
