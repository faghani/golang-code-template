package router

import (
	"strconv"
	"time"

	"github.com/rashadansari/golang-code-template/config"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	LabelEcoCode   = "code"
	LabelEcoMethod = "method"
	LabelEcoHost   = "host"
	LabelEcoURL    = "url"
)

type Metrics struct {
	ReqQPS      *prometheus.CounterVec
	ReqDuration *prometheus.HistogramVec
}

//nolint:gochecknoglobals
var (
	metrics = Metrics{
		ReqQPS: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: config.Namespace,
			Name:      "http_request_total",
		}, []string{LabelEcoCode, LabelEcoMethod, LabelEcoHost, LabelEcoURL}),

		ReqDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Name:      "http_request_duration_seconds",
		}, []string{LabelEcoCode, LabelEcoMethod, LabelEcoHost, LabelEcoURL}),
	}
)

func prometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}

			uri := req.URL.Path
			status := strconv.Itoa(res.Status)
			duration := time.Since(start).Seconds()

			metrics.ReqQPS.WithLabelValues(status, req.Method, req.Host, uri).Inc()
			metrics.ReqDuration.WithLabelValues(status, req.Method, req.Host, uri).Observe(duration)

			return nil
		}
	}
}
