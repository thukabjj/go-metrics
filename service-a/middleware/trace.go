package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func Trace(Tracer trace.Tracer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if span := trace.SpanFromContext(ctx); span == nil {
		}

	}
}
