package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
)

type Pong struct {
	Tracer     trace.Tracer
	HttpClient http.Client
}

type TraceResponse struct {
	Message string `json:"message"`
}

func (h Pong) Pong(ctx *gin.Context) {
	log.WithContext(ctx.Request.Context()).Info("Request")
	ctxReq := ctx.Request.Context()
	requestHeaders := ctx.Request.Header
	ctxChild, span := h.Tracer.Start(ctxReq, "PongHandler")

	trace := h.getTrace(ctxChild, requestHeaders)

	pongResponse := h.getPongFromRepository(ctxChild)

	reponse := pongResponse + trace.Message
	defer span.End()
	ctx.JSON(http.StatusOK, gin.H{"status": reponse})
}

func (h Pong) getTrace(ctx context.Context, headers http.Header) TraceResponse {
	ctxChild, span := h.Tracer.Start(ctx, "getTrace")
	defer span.SetStatus(codes.Ok, "getTrace")
	defer span.End()

	// Make sure you pass the context to the request to avoid broken traces.
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctxChild, carrier)
	//parentCtx := otel.GetTextMapPropagator().Extract(ctxChild, carrier)
	// This carrier is sent accros the process
	fmt.Println(carrier)
	req, err := http.NewRequestWithContext(ctxChild, "GET", "http://java-service-c:8083/trace", nil)

	if err != nil {
		panic(err)
	}

	setHeaderFromPropagatorToRequest(req, carrier)
	// All requests made with this client will create spans.
	res, err := h.HttpClient.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	var traceResponse TraceResponse
	reqBody, err := io.ReadAll(res.Body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	err = json.Unmarshal(reqBody, &traceResponse)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	return traceResponse

}

func setHeaderFromPropagatorToRequest(req *http.Request, carrier propagation.MapCarrier) {
	for _, key := range carrier.Keys() {
		if carrier.Get(key) != "" {
			req.Header.Add(key, carrier.Get(key))
		}
	}
}

func (h Pong) getPongFromRepository(ctx context.Context) string {
	log.WithContext(ctx).Info("getPongFromRepository")
	_, span := h.Tracer.Start(ctx, "getPongFromRepository")
	defer span.End()
	return "pong"
}
