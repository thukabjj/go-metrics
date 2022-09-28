package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/thukabjj/go-metric/service-b/handler"
	"github.com/thukabjj/go-metric/service-b/middleware"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"os"
)

const (
	SERVICE_NAME           = "service-a"
	DEPLOYMENT_ENVIRONMENT = "production"
)

var reqCountMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "go_metrics_microservice_processed_req_total",
	Help: "The total number of processed events",
}, []string{"x_trace_id", "method", "path", "status_code"})

func main() {
	//LOGS
	log.SetFormatter(customLogger{
		formatter: log.JSONFormatter{FieldMap: log.FieldMap{
			"msg": "message",
		}},
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	tp, err := initTracer()
	tracer := tp.Tracer(SERVICE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	prometheus.MustRegister(reqCountMetrics)

	r := gin.Default()
	r.Use(otelgin.Middleware(SERVICE_NAME), middleware.CounterRequestMetrics(reqCountMetrics))

	ping := handler.Ping{
		Tracer: tracer,
	}

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/ping", ping.Ping)
	r.Run()

}

func initTracer() (*sdktrace.TracerProvider, error) {
	// Create the Jaeger exporter
	jaegerHost := getEnv("JAEGER_AGENT_HOST", "http://jaeger:14268/api/traces")
	jaegerExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerHost)))
	if err != nil {
		return nil, err
	}
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(SERVICE_NAME),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String("abcdef12345"),
		semconv.DeploymentEnvironmentKey.String(DEPLOYMENT_ENVIRONMENT),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(jaegerExporter),
		//sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)
	p := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader | b3.B3SingleHeader))

	otel.SetTextMapPropagator(p)
	return tp, nil
}

type customLogger struct {
	formatter log.JSONFormatter
}

func (l customLogger) Format(entry *log.Entry) ([]byte, error) {
	span := trace.SpanFromContext(entry.Context)
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["span_id"] = span.SpanContext().SpanID().String()
	entry.Data["Context"] = span.SpanContext()

	return l.formatter.Format(entry)
}

// getEnv get key environment variable if exist otherwise return defalutValue
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
