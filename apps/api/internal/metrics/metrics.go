package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Metrics
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of HTTP requests to the API Service",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duration of HTTP requests to the API Service",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Business Metrics
	TasksCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "api_tasks_created_total",
			Help: "Total number of tasks created",
		},
	)

	TasksRetrievedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "api_tasks_retrieved_total",
			Help: "Total number of task retrieval requests",
		},
	)

	// Kafka Metrics
	KafkaMessagesPublishedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "api_kafka_messages_published_total",
			Help: "Total number of messages published to Kafka",
		},
	)

	KafkaPublishErrorsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "api_kafka_publish_errors_total",
			Help: "Total number of errors when publishing to Kafka",
		},
	)
)
