# Auth Service - Prometheus Metrics Documentation

## Overview
The Auth Service now exposes Prometheus metrics for monitoring authentication operations, request performance, and business KPIs.

## Metrics Endpoint
```
GET http://localhost:50051/metrics
```

## Available Metrics

### 1. HTTP Request Metrics (RED - Rate, Errors, Duration)

#### `auth_requests_total`
**Type:** Counter  
**Description:** Total number of HTTP requests  
**Labels:**
- `method` - HTTP method (GET, POST)
- `endpoint` - Request endpoint (/register, /login, /validate)
- `status` - HTTP status code (200, 400, 401, 500)

**Example:**
```prometheus
auth_requests_total{endpoint="/login",method="POST",status="200"} 42
auth_requests_total{endpoint="/register",method="POST",status="400"} 3
```

#### `auth_request_duration_seconds`
**Type:** Histogram  
**Description:** HTTP request latency in seconds  
**Labels:**
- `method` - HTTP method
- `endpoint` - Request endpoint

**Buckets:** 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10

**Example PromQL Queries:**
```promql
# P95 latency for login endpoint
histogram_quantile(0.95, rate(auth_request_duration_seconds_bucket{endpoint="/login"}[5m]))

# Average request duration
rate(auth_request_duration_seconds_sum[5m]) / rate(auth_request_duration_seconds_count[5m])
```

### 2. Business Metrics

#### `auth_registrations_total`
**Type:** Counter  
**Description:** Total number of user registrations  
**Labels:** None

**Example:**
```prometheus
auth_registrations_total 156
```

#### `auth_logins_total`
**Type:** Counter  
**Description:** Total number of login attempts  
**Labels:**
- `status` - success or failure

**Example:**
```prometheus
auth_logins_total{status="success"} 1234
auth_logins_total{status="failure"} 45
```

**PromQL - Login Success Rate:**
```promql
rate(auth_logins_total{status="success"}[5m]) 
/ 
rate(auth_logins_total[5m])
```

#### `auth_token_validations_total`
**Type:** Counter  
**Description:** Total number of token validations  
**Labels:**
- `status` - valid or invalid

**Example:**
```prometheus
auth_token_validations_total{status="valid"} 5678
auth_token_validations_total{status="invalid"} 23
```

### 3. Database Metrics

#### `auth_db_query_duration_seconds`
**Type:** Histogram  
**Description:** Database query latency in seconds  
**Labels:**
- `query_type` - select, insert, update

**Buckets:** 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1

*Note: Currently defined but not yet instrumented in database layer*

## Testing Metrics

### 1. View All Metrics
```bash
curl http://localhost:50051/metrics
```

### 2. Filter Specific Metrics
```bash
# View only auth-specific metrics
curl -s http://localhost:50051/metrics | grep "^auth_"

# View registration metrics
curl -s http://localhost:50051/metrics | grep "auth_registrations"
```

### 3. Generate Test Data
```bash
# Register a user
curl -X POST http://localhost:50051/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# Login
curl -X POST http://localhost:50051/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# Check metrics
curl -s http://localhost:50051/metrics | grep -E "(registrations|logins)"
```

## Grafana Dashboard Queries

### Request Rate
```promql
rate(auth_requests_total[5m])
```

### Error Rate
```promql
rate(auth_requests_total{status=~"4..|5.."}[5m])
```

### P95 Latency
```promql
histogram_quantile(0.95, rate(auth_request_duration_seconds_bucket[5m]))
```

### Login Success Rate
```promql
sum(rate(auth_logins_total{status="success"}[5m])) 
/ 
sum(rate(auth_logins_total[5m])) * 100
```

### Active Registrations (Last Hour)
```promql
increase(auth_registrations_total[1h])
```

## Implementation Details

### Files Created/Modified
- `apps/auth/internal/metrics/metrics.go` - Metric definitions
- `apps/auth/internal/middleware/prometheus.go` - HTTP metrics middleware
- `apps/auth/cmd/main.go` - Integration and business metrics

### Middleware Integration
The Prometheus middleware automatically records:
- Request count with labels
- Request duration histogram
- Applied to all routes

### Business Metrics
Manually instrumented at key points:
- Registration success → `metrics.RegistrationsTotal.Inc()`
- Login success/failure → `metrics.LoginsTotal.WithLabelValues(status).Inc()`
- Token validation → `metrics.TokenValidationsTotal.WithLabelValues(status).Inc()`

## Next Steps
1. Deploy Prometheus server
2. Configure scraping for Auth Service
3. Create Grafana dashboard
4. Set up alerting rules
5. Add similar metrics to API and Worker services
