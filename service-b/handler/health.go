package handler

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
)

type Pong struct {
	Tracer     trace.Tracer
	HttpClient http.Client
}

func (h Pong) Pong(ctx *gin.Context) {
	log.WithContext(ctx.Request.Context()).Info("Request")
	ctxReq := ctx.Request.Context()
	ctxChild, span := h.Tracer.Start(ctxReq, "PongHandler")
	pongResponse := h.getPongFromRepository(ctxChild)
	defer span.End()
	ctx.JSON(http.StatusOK, gin.H{"status": pongResponse})
}

func (h Pong) getPongFromRepository(ctx context.Context) string {
	log.WithContext(ctx).Info("getPongFromRepository")
	_, span := h.Tracer.Start(ctx, "getPongFromRepository")
	defer span.End()
	return "pong"
}
