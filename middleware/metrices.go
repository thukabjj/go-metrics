package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

// reqCountMetrices.With(prometheus.Labels{"method": r.Request.Method, "path": r.Request.RequestURI, "statuscode": r.ResponseWriter).Inc()
func CounterRequestMetrics(reqCountMetrics *prometheus.CounterVec) gin.HandlerFunc {
	return func(r *gin.Context) {
		r.Next()
		labels := prometheus.Labels{"method": r.Request.Method, "path": r.Request.RequestURI, "status_code": strconv.Itoa(r.Writer.Status())}
		reqCountMetrics.With(labels).Inc()
	}
}
