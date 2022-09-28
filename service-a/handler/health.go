package handler

import (
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Ping struct {
	Tracer trace.Tracer
}

func (h Ping) Ping(ctx *gin.Context) {
	log.WithContext(ctx.Request.Context()).Info("Request")
	_, span := h.Tracer.Start(ctx.Request.Context(), "Ping")
	defer span.End()
	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}
