package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Ping struct {
	Tracer     trace.Tracer
	HttpClient http.Client
}
type Pong struct {
	Status string `json:"status"`
}

func (h Ping) Ping(ctx *gin.Context) {
	log.WithContext(ctx.Request.Context()).Info("Request")
	ctxRequest := ctx.Request.Context()
	ctxChild, span := h.Tracer.Start(ctxRequest, "PingHandler")
	defer span.End()
	pong := h.getPong(ctxChild)
	ctx.JSON(http.StatusOK, gin.H{"status": pong.Status})
}

func (h Ping) getPong(ctx context.Context) Pong {
	ctxChild, span := h.Tracer.Start(ctx, "getPong")
	defer span.SetStatus(codes.Ok, "getPong")
	defer span.End()

	// Make sure you pass the context to the request to avoid broken traces.
	req, err := http.NewRequestWithContext(ctxChild, "GET", "http://go-service-b:8081/pong", nil)
	if err != nil {
		panic(err)
	}

	// All requests made with this client will create spans.
	res, err := h.HttpClient.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	var pong Pong
	reqBody, err := io.ReadAll(res.Body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	err = json.Unmarshal(reqBody, &pong)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	return pong

}
