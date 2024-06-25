package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func NewLogger() fiber.Handler {
	var (
		once       sync.Once
		errHandler fiber.ErrorHandler
	)

	return func(c *fiber.Ctx) error {
		var bytesSent int
		if v := c.Response().Header.ContentLength(); v < 0 {
			bytesSent = 0
		} else {
			bytesSent = len(c.Response().Body())
		}

		once.Do(func() {
			errHandler = c.App().Config().ErrorHandler
		})

		start := time.Now()
		err := c.Next()
		if err != nil {
			if err := errHandler(c, err); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}
		stop := time.Now()

		log.Info().
			Str("ip", c.IP()).
			Int("status", c.Response().StatusCode()).
			Str("method", c.Method()).
			Str("url", c.OriginalURL()).
			Dur("latency", stop.Sub(start).Round(time.Millisecond)).
			Int("bytes_received", len(c.Request().Body())).
			Int("bytes_sent", bytesSent).
			Err(err).
			Send()
		return nil
	}
}
