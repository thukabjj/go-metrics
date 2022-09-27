package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thukabjj/go-metric/handler"
	"github.com/thukabjj/go-metric/middleware"
)

var reqCountMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "go_metrics_microservice_processed_req_total",
	Help: "The total number of processed events",
}, []string{"method", "path", "status_code"})

func main() {
	prometheus.MustRegister(reqCountMetrics)

	r := gin.Default()
	r.Use(middleware.CounterRequestMetrics(reqCountMetrics))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/health", handler.HealthCheck)
	r.Run()

}
