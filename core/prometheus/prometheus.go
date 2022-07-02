package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

type Labels prometheus.Labels

type Monitor struct {
	ServiceName string //监控服务的名称

	apiRequestsCounter *prometheus.CounterVec
	requestDuration    *prometheus.HistogramVec
	requestSize        *prometheus.HistogramVec
	responseSize       *prometheus.HistogramVec
}

type option func(m *Monitor)

func APIRequestsCounter(labels Labels) option {
	return func(m *Monitor) {
		m.apiRequestsCounter.With(prometheus.Labels(labels)).Inc()
	}
}

func RequestDuration(labels Labels, duration float64) option {
	return func(m *Monitor) {
		m.requestDuration.With(prometheus.Labels(labels)).Observe(duration)
	}
}

func RequestSize(labels Labels, reqSize float64) option {
	return func(m *Monitor) {
		m.requestSize.With(prometheus.Labels(labels)).Observe(reqSize)
	}
}

func ResponseSize(labels Labels, respSize float64) option {
	return func(m *Monitor) {
		m.responseSize.With(prometheus.Labels(labels)).Observe(respSize)
	}
}

func NewMonitor(namespace string) *Monitor {
	namespace = strings.ReplaceAll(namespace, "-", "_")

	apiRequestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_requests_total",
			Help:      "A counter for requests to the wrapped handler.",
		},
		[]string{"handler", "method", "http_code", "code", "err_msg"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "A histogram of latencies for requests.",
		},
		[]string{"handler", "method", "http_code", "code"},
	)

	requestSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_size_bytes",
			Help:      "A histogram of request sizes for requests.",
		},
		[]string{"handler", "method", "http_code", "code"},
	)

	responseSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_response_size_bytes",
			Help:      "A histogram of response sizes for requests.",
		},
		[]string{"handler", "method", "http_code", "code"},
	)

	//注册指标
	prometheus.MustRegister(apiRequestsCounter, requestDuration, requestSize, responseSize)

	return &Monitor{
		apiRequestsCounter: apiRequestsCounter,
		requestDuration:    requestDuration,
		requestSize:        requestSize,
		responseSize:       responseSize,
	}
}

func (m *Monitor) With(opts ...option) {
	for _, opt := range opts {
		opt(m)
	}
}
