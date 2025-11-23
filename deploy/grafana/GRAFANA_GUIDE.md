# Grafana Setup & Usage Guide

## Access Grafana

**Grafana UI:** http://localhost:3001  
**Login:** `admin` / `admin`

## Pre-Configured Components

### Datasource
- **Prometheus** - Auto-configured to connect to Prometheus server
- URL: http://prometheus:9090
- Scrape interval: 15 seconds

### Dashboards
- **Auth Service - Overview** - Pre-loaded dashboard with 7 panels

## Auth Service Dashboard

The Auth Service dashboard includes the following panels:

### 1. Request Rate (req/sec)
Shows requests per second for all endpoints
- **Query:** `rate(auth_requests_total[5m])`
- **Visualization:** Time series graph
- **Legend:** Endpoint and HTTP method

### 2. P95 Latency
Shows 95th percentile and median latency
- **Queries:** 
  - P95: `histogram_quantile(0.95, rate(auth_request_duration_seconds_bucket[5m]))`
  - P50: `histogram_quantile(0.50, rate(auth_request_duration_seconds_bucket[5m]))`
- **Unit:** Seconds

### 3. Total Registrations
Current count of user registrations
- **Query:** `auth_registrations_total`
- **Visualization:** Stat panel

### 4. Login Success Rate (%)
Percentage of successful logins
- **Query:** `sum(rate(auth_logins_total{status="success"}[5m])) / sum(rate(auth_logins_total[5m])) * 100`
- **Range:** 0-100%

### 5. Token Validations
Rate of valid vs invalid token validations
- **Queries:**
  - Valid: `rate(auth_token_validations_total{status="valid"}[5m])`
  - Invalid: `rate(auth_token_validations_total{status="invalid"}[5m])`

### 6. Login Attempts
Success vs failure login attempts
- **Queries:**
  - Success: `rate(auth_logins_total{status="success"}[5m])`
  - Failure: `rate(auth_logins_total{status="failure"}[5m])`

### 7. Error Rate
HTTP 4xx and 5xx errors
- **Query:** `rate(auth_requests_total{status=~"4..|5.."}[5m])`
- **Alert:** Triggers if error rate > 5 req/sec

## Using the Dashboard

### View Dashboard
1. Open http://localhost:3001
2. Login with admin/admin
3. Go to "Dashboards" → "Auth Service - Overview"

### Customize Time Range
- Top right corner: Select time range (Last 1h, 6h, 24h, etc.)
- Refresh interval: 10 seconds (auto-refresh)

### Zoom In
- Click and drag on any graph to zoom into a specific time range
- Click "Zoom out" to reset

### View Panel Details
- Click on any panel title → "Edit" to see the query
- Click "View" to see full-screen visualization

## Creating Additional Dashboards

### For API Service (Future)
```json
{
  "title": "API Service - Overview",
  "panels": [
    {
      "title": "Task Creation Rate",
      "expr": "rate(api_tasks_created_total[5m])"
    },
    {
      "title": "Kafka Publish Errors",
      "expr": "rate(api_kafka_publish_errors_total[5m])"
    }
  ]
}
```

### For Worker Service (Future)
```json
{
  "title": "Worker Service - Overview",
  "panels": [
    {
      "title": "Tasks Processed",
      "expr": "rate(worker_tasks_processed_total[5m])"
    },
    {
      "title": "Kafka Consumer Lag",
      "expr": "worker_kafka_lag"
    }
  ]
}
```

## Troubleshooting

### Dashboard Not Showing
```bash
# Check if Grafana is running
docker ps | grep pixelflow-grafana

# Check Grafana logs
docker logs pixelflow-grafana

# Restart Grafana
docker restart pixelflow-grafana
```

### No Data in Panels
1. Check Prometheus datasource: Configuration → Data Sources → Prometheus
2. Click "Test" - should show "Data source is working"
3. Verify Prometheus has data: http://localhost:9091
4. Check time range (top right) - try "Last 1 hour"

### Datasource Connection Failed
```bash
# Verify Prometheus is running
docker ps | grep pixelflow-prometheus

# Check if services are on same network
docker network inspect harmonic-rosette_pixelflow-net

# Restart both services
docker restart pixelflow-prometheus pixelflow-grafana
```

## Advanced Features

### Alerting
The Error Rate panel has a pre-configured alert:
- **Condition:** Error rate > 5 req/sec
- **Frequency:** Check every 1 minute
- **Action:** Currently set to log only

To configure notifications:
1. Go to Alerting → Notification channels
2. Add email, Slack, PagerDuty, etc.
3. Edit panel alert to use notification channel

### Variables
Add dashboard variables for dynamic filtering:
1. Dashboard settings → Variables → Add variable
2. Example: `endpoint` variable to filter by endpoint
3. Query: `label_values(auth_requests_total, endpoint)`

### Annotations
Mark important events on graphs:
1. Dashboard settings → Annotations
2. Add annotation query
3. Example: Mark deployments, incidents, etc.

## Exporting Dashboards

### Export as JSON
1. Dashboard settings → JSON Model
2. Copy JSON
3. Save to file

### Import Dashboard
1. Dashboards → Import
2. Upload JSON file or paste JSON
3. Select Prometheus datasource

## Useful Links

- Grafana UI: http://localhost:3001
- Prometheus: http://localhost:9091
- Auth Metrics: http://localhost:50051/metrics
- Grafana Docs: https://grafana.com/docs/
