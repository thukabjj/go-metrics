package handler

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Pong struct {
	Tracer trace.Tracer
}

func (h Pong) Pong(ctx *gin.Context) {
	log.WithContext(ctx.Request.Context()).Info("Request")
	_, span := h.Tracer.Start(ctx.Request.Context(), "Pong")
	pongResponse := h.getPongFromRepository(ctx.Request.Context())
	defer span.End()
	ctx.JSON(http.StatusOK, gin.H{"status": pongResponse})
}

func (h Pong) getPongFromRepository(ctx context.Context) string {
	log.WithContext(ctx).Info("getPongFromRepository")
	_, span := h.Tracer.Start(ctx, "getPongFromRepository")
	defer span.End()
	return "pong"
}
