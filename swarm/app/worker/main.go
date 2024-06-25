package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"example.com/shared/redis"
	"example.com/shared/telemetry"
	"example.com/worker/constant"
	"example.com/worker/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if os.Getenv("DEBUG") != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	}

	app := fiber.New(fiber.Config{
		AppName:               constant.AppID,
		DisableStartupMessage: os.Getenv("APP_ENV") == "production",
	})
	app.Use(pprof.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).Send([]byte("ok"))
	})

	telemetry := telemetry.New(ctx, "example.com/worker", constant.AppID)
	redis := redis.New(ctx, constant.AppID)
	handler := handler.New(telemetry.Tracer)

	pubsub := redis.PSubscribe(ctx, "user.*")

	app.Hooks().OnShutdown(func() error {
		pubsub.Close()
		redis.Close()
		return nil
	})

	go func() {
		for msg := range pubsub.Channel() {
			handler.Process(ctx, msg.Channel, msg.Payload)
		}
	}()

	go func() {
		if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	stop()
	app.Shutdown()
}
