package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	tracer trace.Tracer
}

func New(tracer trace.Tracer) *Handler {
	return &Handler{tracer}
}

func (h *Handler) Process(ctx context.Context, channel, v string) {
	var payload struct {
		Body       map[string]interface{} `json:"body"`
		Attributes map[string]string      `json:"attributes"`
	}
	if err := json.Unmarshal([]byte(v), &payload); err != nil {
		log.Error().Err(err).Msg("handler.Receive")
		return
	}

	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(payload.Attributes))

	attrs := []attribute.KeyValue{
		semconv.MessagingOperationDeliver,
		semconv.MessagingDestinationName(channel),
	}
	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	}
	ctx, span := h.tracer.Start(ctx, fmt.Sprintf("%s process", channel), opts...)
	defer span.End()

	switch channel {
	case "user.poke":
		h.Poke(ctx, payload.Body)
	}
}
