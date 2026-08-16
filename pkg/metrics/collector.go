package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ProbeSuccess = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "healthgate_probe_success",
			Help: "1 if probe succeeded, 0 otherwise",
		},
		[]string{"type", "target"},
	)

	ProbeLatencySeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "healthgate_probe_duration_seconds",
			Help:    "Probe round-trip duration in seconds",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
		},
		[]string{"type", "target"},
	)
)
