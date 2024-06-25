package handler

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (h *Handler) Poke(ctx context.Context, body map[string]interface{}) {
	if v, ok := body["userId"].(string); !ok {
		return
	} else {
		log.Info().Str("userId", v).Msg("poked")
	}
}
