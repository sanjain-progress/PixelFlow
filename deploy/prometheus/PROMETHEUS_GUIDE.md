# Prometheus Setup & Usage Guide

## Overview
Prometheus is now deployed and actively scraping metrics from the Auth Service. This guide shows you how to use Prometheus to monitor your application.

## Access Prometheus

**Prometheus UI:** http://localhost:9091

## Current Status

### Active Targets
- ✅ **prometheus** - Prometheus self-monitoring (UP)
- ✅ **auth-service** - Auth Service metrics (UP)
- ⏸️ **api-service** - Not yet instrumented (DOWN)
- ⏸️ **worker-service** - Not yet instrumented (DOWN)

### Metrics Being Collected
Auth Service is currently exposing:
- `auth_requests_total` - 1051 validation requests, 1 registration, 1 login
- `auth_request_duration_seconds` - Request latency histogram
- `auth_registrations_total` - User registrations
- `auth_logins_total` - Login attempts
- `auth_token_validations_total` - Token validations

## Using Prometheus

### 1. View Targets
Check which services Prometheus is scraping:
```
http://localhost:9091/targets
```

### 2. Execute Queries

#### Basic Queries
```promql
# Total requests to Auth Service
auth_requests_total

# Requests by endpoint
auth_requests_total{endpoint="/login"}

# Successful logins
auth_logins_total{status="success"}

# Total registrations
auth_registrations_total
```

#### Rate Queries (Requests per second)
```promql
# Request rate for last 5 minutes
rate(auth_requests_total[5m])

# Login rate
rate(auth_logins_total[5m])
```

#### Latency Queries
```promql
# P95 latency for all endpoints
histogram_quantile(0.95, rate(auth_request_duration_seconds_bucket[5m]))

# P95 latency for login endpoint
histogram_quantile(0.95, rate(auth_request_duration_seconds_bucket{endpoint="/login"}[5m]))

# Average latency
rate(auth_request_duration_seconds_sum[5m]) / rate(auth_request_duration_seconds_count[5m])
```

#### Business Metrics
```promql
# Login success rate (percentage)
sum(rate(auth_logins_total{status="success"}[5m])) 
/ 
sum(rate(auth_logins_total[5m])) * 100

# Registrations in last hour
increase(auth_registrations_total[1h])

# Token validation success rate
sum(rate(auth_token_validations_total{status="valid"}[5m])) 
/ 
sum(rate(auth_token_validations_total[5m])) * 100
```

### 3. Graph Metrics

1. Go to http://localhost:9091/graph
2. Enter a PromQL query
3. Click "Execute"
4. Switch to "Graph" tab to visualize

**Example Graphs:**
- Request rate over time: `rate(auth_requests_total[5m])`
- Login attempts: `rate(auth_logins_total[5m])`
- P95 latency: `histogram_quantile(0.95, rate(auth_request_duration_seconds_bucket[5m]))`

## Testing Metrics

### Generate Traffic
```bash
# Register a user
curl -X POST http://localhost:50051/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# Login
curl -X POST http://localhost:50051/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# Wait 15 seconds for Prometheus to scrape
sleep 15

# Query metrics
curl -s 'http://localhost:9091/api/v1/query?query=auth_registrations_total' | jq
```

### Verify Scraping
```bash
# Check if Auth Service is UP
curl -s 'http://localhost:9091/api/v1/targets' | jq '.data.activeTargets[] | select(.labels.job=="auth-service")'

# Query auth metrics
curl -s 'http://localhost:9091/api/v1/query?query=auth_requests_total' | jq
```

## Prometheus Configuration

### Scrape Configuration
Located at: `deploy/prometheus/prometheus.yml`

```yaml
scrape_configs:
  - job_name: 'auth-service'
    static_configs:
      - targets: ['pixelflow-auth:50051']
        labels:
          service: 'auth'
          app: 'pixelflow'
```

### Scrape Interval
- **Default:** 15 seconds
- Prometheus scrapes `/metrics` endpoint every 15s

### Data Retention
- **Default:** 15 days
- Stored in Docker volume: `prometheus_data`

## Common PromQL Patterns

### Counters (always increasing)
```promql
# Use rate() for per-second rate
rate(auth_requests_total[5m])

# Use increase() for total increase over time
increase(auth_requests_total[1h])
```

### Histograms (latency)
```promql
# Quantiles (P50, P95, P99)
histogram_quantile(0.50, rate(auth_request_duration_seconds_bucket[5m]))
histogram_quantile(0.95, rate(auth_request_duration_seconds_bucket[5m]))
histogram_quantile(0.99, rate(auth_request_duration_seconds_bucket[5m]))
```

### Aggregation
```promql
# Sum across all endpoints
sum(rate(auth_requests_total[5m]))

# Group by endpoint
sum by (endpoint) (rate(auth_requests_total[5m]))

# Count unique endpoints
count(auth_requests_total)
```

## Troubleshooting

### Target is DOWN
```bash
# Check if service is running
docker ps | grep pixelflow-auth

# Check if /metrics endpoint is accessible
curl http://localhost:50051/metrics

# Check Prometheus logs
docker logs pixelflow-prometheus
```

### No Data Showing
- Wait 15 seconds for first scrape
- Verify target is UP in http://localhost:9091/targets
- Check if metrics exist: `curl http://localhost:50051/metrics`

### Metrics Not Updating
- Generate traffic to the service
- Wait for next scrape interval (15s)
- Verify time range in Prometheus UI

## Next Steps

1. **Add Grafana** - Visualize metrics with dashboards
2. **Instrument API Service** - Add /metrics endpoint
3. **Instrument Worker Service** - Add /metrics endpoint
4. **Create Alerts** - Set up AlertManager
5. **Add Recording Rules** - Pre-compute expensive queries

## Useful Links

- Prometheus UI: http://localhost:9091
- Targets: http://localhost:9091/targets
- Graph: http://localhost:9091/graph
- Auth Metrics: http://localhost:50051/metrics

## API Endpoints

### Prometheus HTTP API
```bash
# Query instant value
curl 'http://localhost:9091/api/v1/query?query=auth_requests_total'

# Query range
curl 'http://localhost:9091/api/v1/query_range?query=rate(auth_requests_total[5m])&start=2024-01-01T00:00:00Z&end=2024-01-01T01:00:00Z&step=15s'

# Get targets
curl 'http://localhost:9091/api/v1/targets'

# Get labels
curl 'http://localhost:9091/api/v1/labels'
```
