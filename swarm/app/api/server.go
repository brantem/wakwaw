package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"example.com/api/constant"
	"example.com/api/handler"
	"example.com/api/middleware"
	"example.com/shared/db"
	"example.com/shared/redis"
	"example.com/shared/telemetry"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	isDebug := os.Getenv("DEBUG") == "1"
	if isDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	}

	telemetry := telemetry.New(ctx, "example.com/api", constant.AppID)
	db := db.New(constant.AppID)
	redis := redis.New(ctx, constant.AppID)
	app := fiber.New(fiber.Config{
		AppName:               constant.AppID,
		DisableStartupMessage: os.Getenv("APP_ENV") == "production",
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Hooks().OnShutdown(func() error {
		db.Close()
		redis.Close()
		telemetry.Shutdown(ctx)
		return nil
	})

	app.Use(otelfiber.Middleware())
	// app.Use(pprof.New())

	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	app.Use(etag.New())
	app.Use(helmet.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: isDebug}))
	app.Use(middleware.NewLogger())

	h := handler.New(telemetry.Tracer, db, redis)
	h.Register(app)

	go func() {
		if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	stop()
	app.Shutdown()
}
