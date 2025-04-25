package web

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)
	serverErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "server_errors_total",
			Help: "Total number of server errors by route and status code",
		},
		[]string{"route", "status_code"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handler",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func requestsMetric(c *gin.Context) {
	start := time.Now()

	c.Next()

	duration := time.Since(start).Seconds()
	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}

	httpRequestsTotal.WithLabelValues(c.Request.Method, path).Inc()
	httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
}

func serverErrorsMetric(c *gin.Context) {
	c.Next()

	status := c.Writer.Status()
	route := c.FullPath()

	if route == "" {
		route = "unknown"
	}

	if status >= 400 {
		serverErrors.WithLabelValues(route, strconv.Itoa(status)).Inc()
	}
}
