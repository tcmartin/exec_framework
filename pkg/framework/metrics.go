package framework

import (
    "github.com/prometheus/client_golang/prometheus"
)

// Metrics holds Prometheus collectors
type Metrics struct {
    NodeDuration *prometheus.HistogramVec
    NodeErrors   *prometheus.CounterVec
}

// NewMetrics registers and returns metrics
func NewMetrics() *Metrics {
    m := &Metrics{
        NodeDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{Namespace: "workflow", Name: "node_duration_seconds"},
            []string{"node"},
        ),
        NodeErrors: prometheus.NewCounterVec(
            prometheus.CounterOpts{Namespace: "workflow", Name: "node_errors_total"},
            []string{"node"},
        ),
    }
    prometheus.MustRegister(m.NodeDuration, m.NodeErrors)
    return m
}
