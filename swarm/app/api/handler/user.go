package handler

import (
	"database/sql"
	"encoding/json"

	"example.com/api/constant"
	"example.com/api/model"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.23.0"
	"go.opentelemetry.io/otel/trace"
)

func (h *Handler) GetUser(c *fiber.Ctx) error {
	var result struct {
		User  *model.User `json:"user"`
		Error any         `json:"error"`
	}

	userID := c.Params("userId")

	span := trace.SpanFromContext(c.UserContext())
	span.SetAttributes(attribute.String("user.id", userID))

	var user model.User
	err := h._dbRead.QueryRowxContext(c.UserContext(), `
		SELECT id, name
		FROM users
		WHERE id = $1
	`, userID).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			result.Error = constant.RespNotFound
			return c.Status(fiber.StatusNotFound).JSON(result)
		}
		log.Error().Err(err).Msg("user.GetUser")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	log.Debug().Any("user", user).Msg("user.GetUser")
	result.User = &user

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) PokeUser(c *fiber.Ctx) error {
	userID := c.Params("userId")

	attrs := []attribute.KeyValue{
		attribute.String("user.id", userID),
		semconv.MessagingOperationPublish,
		semconv.MessagingDestinationName("user.poke"),
	}
	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}
	ctx, span := h.tracer.Start(c.UserContext(), "user.poke publish", opts...)
	defer span.End()

	attributes := &propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, attributes)

	m := map[string]interface{}{
		"body":       map[string]string{"userId": userID},
		"attributes": attributes,
	}
	buf, _ := json.Marshal(m)

	h.redis.Publish(ctx, "user.poke", string(buf))
	return c.SendStatus(fiber.StatusAccepted)
}
