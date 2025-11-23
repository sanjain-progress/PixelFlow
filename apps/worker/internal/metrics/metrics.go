package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Task Processing Metrics
	TasksProcessedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_tasks_processed_total",
			Help: "Total number of tasks processed by the worker",
		},
		[]string{"status"}, // success, failure
	)

	TaskProcessingDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "worker_task_processing_duration_seconds",
			Help:    "Duration of task processing",
			Buckets: prometheus.DefBuckets,
		},
	)

	ActiveProcessingTasks = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_active_processing_tasks",
			Help: "Number of tasks currently being processed",
		},
	)

	// Kafka Metrics
	KafkaMessagesConsumedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "worker_kafka_messages_consumed_total",
			Help: "Total number of messages consumed from Kafka",
		},
	)

	KafkaConsumptionErrorsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "worker_kafka_consumption_errors_total",
			Help: "Total number of errors when consuming from Kafka",
		},
	)
)
