package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Request metrics (RED - Rate, Errors, Duration)
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets, // 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
		},
		[]string{"method", "endpoint"},
	)

	// Business metrics
	RegistrationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	LoginsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_logins_total",
			Help: "Total number of login attempts",
		},
		[]string{"status"}, // success, failure
	)

	TokenValidationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_token_validations_total",
			Help: "Total number of token validations",
		},
		[]string{"status"}, // valid, invalid
	)

	// Database metrics
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_db_query_duration_seconds",
			Help:    "Database query latency in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"query_type"}, // select, insert, update
	)
)
