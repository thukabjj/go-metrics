package main

import (
	"context"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
	"github.com/thukabjj/go-metric/service-a/handler"
	"github.com/thukabjj/go-metric/service-a/middleware"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	SERVICE_NAME           = "service-b"
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

	r := gin.New()

	r.Use(otelgin.Middleware(SERVICE_NAME), middleware.CounterRequestMetrics(reqCountMetrics))

	pong := handler.Pong{
		Tracer: tracer,
	}

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/pong", pong.Pong)
	r.Run(":8081")

}

func initTracer() (*sdktrace.TracerProvider, error) {
	log.Println("Initialising tracer")
	log.Println("Connecting to GRPC endpoint...")

	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("otel-collector:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))

	if err != nil {
		return nil, err
	}

	log.Println("Connection established.")

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(SERVICE_NAME),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String("abcdef12345"),
		semconv.DeploymentEnvironmentKey.String(DEPLOYMENT_ENVIRONMENT),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
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
