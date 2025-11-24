# PixelFlow Observability User Guide

This guide provides step-by-step instructions for using the observability tools in the PixelFlow application.

---

## Table of Contents
1. [Grafana Dashboards](#grafana-dashboards)
2. [Prometheus Metrics](#prometheus-metrics)
3. [Jaeger Distributed Tracing](#jaeger-distributed-tracing)
4. [Loki Log Aggregation](#loki-log-aggregation)
5. [Grafana Alerting](#grafana-alerting)

---

## Grafana Dashboards

Grafana provides visualization dashboards for monitoring service metrics.

### Accessing Grafana

1. **Open Grafana UI**
   - Navigate to: `http://localhost:3001`
   - Default credentials:
     - Username: `admin`
     - Password: `admin`

2. **First Login**
   - You'll be prompted to change the password
   - Click "Skip" if you're in a development environment
   - For production, set a strong password

### Viewing Dashboards

1. **Navigate to Dashboards**
   - Click the "Dashboards" icon (four squares) in the left sidebar
   - Or go directly to: `http://localhost:3001/dashboards`

2. **Available Dashboards**
   - **Auth Service Dashboard**: Metrics for authentication service
   - **API Service Dashboard**: Metrics for API service
   - **Worker Service Dashboard**: Metrics for background worker

3. **Opening a Dashboard**
   - Click on any dashboard name to open it
   - Example: Click "API Service Dashboard"

### Understanding Dashboard Panels

Each dashboard contains multiple panels showing different metrics:

**Auth Service Dashboard:**
- Request Rate (requests/second)
- Request Duration (P50, P95, P99 latencies)
- Error Rate (percentage of failed requests)
- User Registrations (total count)
- Active Logins (total count)

**API Service Dashboard:**
- HTTP Request Rate
- Request Duration by Endpoint
- Error Rate by Status Code
- Tasks Created (business metric)
- Kafka Publish Errors

**Worker Service Dashboard:**
- Task Processing Rate
- Task Processing Duration
- Active Processing Tasks (current workload)
- Kafka Messages Consumed
- Task Success/Failure Rate

### Customizing Time Range

1. Click the time picker in the top-right corner
2. Select a preset range (Last 5 minutes, Last 1 hour, etc.)
3. Or set a custom range with "Absolute time range"
4. Click "Apply time range"

### Refreshing Data

- Click the refresh icon (circular arrow) in the top-right
- Or set auto-refresh interval from the dropdown next to it
- Options: 5s, 10s, 30s, 1m, 5m, etc.

---

## Prometheus Metrics

Prometheus collects and stores time-series metrics from all services.

### Accessing Prometheus UI

1. **Open Prometheus**
   - Navigate to: `http://localhost:9091`
   - No authentication required

2. **Main Interface**
   - You'll see the query interface with a text box

### Querying Metrics

1. **View Available Metrics**
   - Click the "Metrics Explorer" icon (globe) next to the query box
   - Browse through all available metrics
   - Metrics are prefixed by service: `auth_`, `api_`, `worker_`

2. **Execute a Query**
   - Type a metric name in the query box
   - Example: `api_requests_total`
   - Click "Execute" or press Enter
   - View results in "Table" or "Graph" tab

3. **Example Queries**

   **Request Rate (last 5 minutes):**
   ```promql
   rate(api_requests_total[5m])
   ```

   **Error Rate:**
   ```promql
   rate(api_requests_total{status=~"5.."}[5m]) / rate(api_requests_total[5m]) * 100
   ```

   **P95 Latency:**
   ```promql
   histogram_quantile(0.95, rate(api_request_duration_seconds_bucket[5m]))
   ```

   **Active Worker Tasks:**
   ```promql
   worker_active_processing_tasks
   ```

### Viewing Service Targets

1. Click "Status" in the top menu
2. Select "Targets"
3. View all scraped services and their health status
4. Services should show "UP" in green

### Checking Service Discovery

1. Click "Status" → "Service Discovery"
2. View all discovered services
3. Verify `auth-service`, `api-service`, and `worker-service` are listed

---

## Jaeger Distributed Tracing

Jaeger visualizes request flows across multiple services.

### Accessing Jaeger UI

1. **Open Jaeger**
   - Navigate to: `http://localhost:16686`
   - No authentication required

### Finding Traces

1. **Search for Traces**
   - In the left sidebar, you'll see the search interface
   - **Service**: Select a service from the dropdown
     - Options: `auth-service`, `api-service`, `worker-service`
   - **Operation**: Select an operation (e.g., `POST /tasks`, `POST /register`)
   - **Lookback**: Set time range (default: Last 1 hour)
   - Click "Find Traces" button

2. **Example: View API Request Trace**
   - Service: `api-service`
   - Operation: `POST /tasks`
   - Click "Find Traces"
   - You'll see a list of traces with duration and span count

### Viewing Trace Details

1. **Click on a Trace**
   - Click any trace from the results list
   - You'll see the trace timeline view

2. **Understanding the Timeline**
   - **Horizontal bars**: Each bar represents a span (operation)
   - **Length**: Duration of the operation
   - **Hierarchy**: Nested spans show parent-child relationships
   - **Colors**: Different services have different colors

3. **Trace Flow Example**
   ```
   api-service: POST /tasks (100ms)
   ├─ api-service: validate_token (10ms)
   ├─ api-service: save_to_mongodb (20ms)
   └─ api-service: publish_to_kafka (15ms)
       └─ worker-service: process_task (50ms)
           ├─ worker-service: fetch_image (30ms)
           └─ worker-service: update_mongodb (10ms)
   ```

4. **Inspecting Span Details**
   - Click on any span to expand it
   - View tags (metadata like `http.method`, `http.status_code`)
   - View logs (events that occurred during the span)
   - View process information (service name, version)

### Analyzing Performance

1. **Identify Slow Operations**
   - Look for long horizontal bars
   - Check which service/operation is taking the most time

2. **Compare Traces**
   - Search for the same operation multiple times
   - Compare durations to identify anomalies
   - Use "Compare" feature to view side-by-side

### Trace Context Propagation

Traces flow through:
1. **HTTP Request** → API Service (creates trace)
2. **Kafka Message** → Worker Service (continues trace)
3. All operations share the same **Trace ID**

---

## Loki Log Aggregation

Loki centralizes logs from all Docker containers.

### Accessing Logs in Grafana

1. **Navigate to Explore**
   - Click the "Explore" icon (compass) in the left sidebar
   - Or go to: `http://localhost:3001/explore`

2. **Select Loki Datasource**
   - In the top-left dropdown, select "Loki"
   - The query builder will appear

### Querying Logs

1. **Switch to Code Mode**
   - Click the "Code" button (if in Builder mode)
   - You'll see a text editor for LogQL queries

2. **Basic Log Query**
   ```logql
   {container_name="pixelflow-api"}
   ```
   - This shows all logs from the API service
   - Click "Run query" (blue button)

3. **Filter by Service**
   ```logql
   {container_name="pixelflow-auth"}     # Auth service logs
   {container_name="pixelflow-worker"}   # Worker service logs
   {container_name="pixelflow-kafka"}    # Kafka logs
   ```

4. **Filter by Log Level**
   ```logql
   {container_name="pixelflow-api"} |= "error"
   {container_name="pixelflow-api"} |= "ERROR"
   ```

5. **Search for Specific Text**
   ```logql
   {container_name="pixelflow-api"} |= "task_id"
   {container_name="pixelflow-worker"} |= "Processing task"
   ```

6. **Time Range Filtering**
   ```logql
   {container_name="pixelflow-api"} | json | level="error"
   ```

### Viewing Log Details

1. **Expand Log Lines**
   - Click on any log line to expand it
   - View full JSON structure (if using structured logging)
   - See timestamp, level, message, and metadata

2. **Copy Log Content**
   - Click the copy icon next to any log line
   - Useful for sharing or debugging

### Live Tailing

1. Click the "Live" button in the top-right
2. Logs will stream in real-time
3. Click "Stop" to pause streaming

### Correlating Logs with Traces

1. **Find Trace ID in Logs**
   - Structured logs include `trace_id` field
   - Copy the trace ID

2. **Search in Jaeger**
   - Go to Jaeger UI
   - Paste trace ID in the search box
   - View the corresponding trace

---

## Grafana Alerting

Grafana alerts notify you when metrics exceed thresholds.

### Viewing Alert Rules

1. **Navigate to Alerting**
   - Click the "Alerting" icon (bell) in the left sidebar
   - Or go to: `http://localhost:3001/alerting/list`

2. **View Alert Rules**
   - You'll see the "Alert rules" page
   - Expand "PixelFlow Alerts" group
   - View all configured alert rules:
     - High Error Rate - API Service
     - High Latency - API Service
     - Service Down - API Service
     - High Error Rate - Worker Service
     - High Kafka Consumption Errors

### Understanding Alert States

**Alert States:**
- 🟢 **Normal**: Metric is within threshold
- 🟡 **Pending**: Threshold exceeded, waiting for evaluation period
- 🔴 **Firing**: Alert is active and should trigger notifications
- ⚪ **No Data**: No metrics received

### Viewing Alert Details

1. **Click on an Alert Rule**
   - Example: Click "High Error Rate - API Service"
   - You'll see the rule configuration page

2. **Rule Details**
   - **Query**: The PromQL query being evaluated
   - **Condition**: Threshold that triggers the alert
   - **Evaluation**: How often the rule is checked (1m)
   - **For**: How long condition must be true before firing (5m)

3. **View Alert History**
   - Scroll down to see "State history"
   - View when alerts fired and resolved
   - See evaluation results over time

### Testing Alerts

#### Method 1: Stop a Service

1. **Stop API Service**
   ```bash
   docker-compose stop api-service
   ```

2. **Wait 2-3 Minutes**
   - The "Service Down - API Service" alert evaluates every 1 minute
   - It fires after 2 minutes of no metrics

3. **Check Alert Status**
   - Go to: `http://localhost:3001/alerting/list`
   - Look for "Service Down - API Service"
   - It should show 🔴 **Firing** status

4. **View Alert Details**
   - Click on the firing alert
   - See the evaluation graph showing the service is down
   - View the alert message

5. **Restart Service**
   ```bash
   docker-compose start api-service
   ```
   - Alert should resolve after ~2 minutes

#### Method 2: Generate High Error Rate

1. **Send Invalid Requests**
   ```bash
   # Send requests that will fail
   for i in {1..100}; do
     curl -X POST http://localhost:8080/tasks \
       -H "Authorization: Bearer invalid_token" \
       -H "Content-Type: application/json" \
       -d '{"image_url":"test.jpg"}'
   done
   ```

2. **Wait 5-6 Minutes**
   - Error rate alert has a 5-minute evaluation window

3. **Check Alert**
   - "High Error Rate - API Service" should fire
   - View in Grafana alerting page

### Configuring Notifications

1. **Navigate to Contact Points**
   - Go to: `http://localhost:3001/alerting/notifications`
   - Click "Contact points" tab

2. **Add Contact Point**
   - Click "New contact point"
   - **Name**: e.g., "Email Notifications"
   - **Type**: Select integration (Email, Slack, PagerDuty, etc.)
   - Configure integration-specific settings
   - Click "Save contact point"

3. **Create Notification Policy**
   - Go to "Notification policies" tab
   - Click "New policy"
   - **Match labels**: e.g., `severity=critical`
   - **Contact point**: Select your contact point
   - Click "Save policy"

4. **Test Notification**
   - Click "Test" button next to your contact point
   - Verify you receive the test notification

### Silencing Alerts

1. **Create Silence**
   - Go to: `http://localhost:3001/alerting/silences`
   - Click "New silence"
   - **Matchers**: Select alert to silence
   - **Duration**: How long to silence
   - **Comment**: Reason for silencing
   - Click "Create"

2. **View Active Silences**
   - See all active silences
   - Expire or delete silences as needed

---

## Quick Reference

### Service URLs
- **Grafana**: http://localhost:3001 (admin/admin)
- **Prometheus**: http://localhost:9091
- **Jaeger**: http://localhost:16686
- **Loki**: Accessible via Grafana Explore

### Common Tasks

**Check Service Health:**
```bash
docker-compose ps
```

**View Service Logs:**
```bash
docker logs -f pixelflow-api
docker logs -f pixelflow-worker
docker logs -f pixelflow-auth
```

**Restart Observability Stack:**
```bash
docker-compose restart prometheus grafana jaeger loki promtail
```

**Generate Test Traffic:**
```bash
./test_e2e.sh
```

### Troubleshooting

**Grafana not loading:**
- Check: `docker logs pixelflow-grafana`
- Restart: `docker-compose restart grafana`

**No metrics in Prometheus:**
- Check targets: http://localhost:9091/targets
- Verify services are exposing `/metrics` endpoints

**No traces in Jaeger:**
- Verify services are running
- Check Jaeger logs: `docker logs pixelflow-jaeger`
- Ensure OTLP port 4317 is accessible

**No logs in Loki:**
- Check Promtail: `docker logs pixelflow-promtail`
- Verify Docker socket is mounted: `/var/run/docker.sock`
