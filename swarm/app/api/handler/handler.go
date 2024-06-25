package handler

import (
	"example.com/api/constant"
	"example.com/api/handler/internal/storage"
	"example.com/shared/db"
	"example.com/shared/redis"
	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	tracer trace.Tracer

	redis   redis.RedisInterface
	_db     *sqlx.DB
	db      sq.StatementBuilderType
	_dbRead *sqlx.DB
	dbRead  sq.StatementBuilderType
}

func New(tracer trace.Tracer, db *db.DB, redis redis.RedisInterface) *Handler {
	h := Handler{tracer: tracer, redis: redis}
	if db != nil {
		sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		h._db = db.Master
		h.db = sb.RunWith(db.Master)
		h._dbRead = db.Read
		h.dbRead = sb.RunWith(db.Read)
	}
	return &h
}

func (h *Handler) Register(r *fiber.App) {
	storage := storage.New(h.redis, constant.StorageKeyPrefix)
	c := cache.New(cache.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.OriginalURL())
		},
		Storage:              storage,
		StoreResponseHeaders: true,
	})

	v1 := r.Group("/v1")

	userID := v1.Group("/users/:userId")
	userID.Get("/", c, h.GetUser)
	userID.Post("/poke", h.PokeUser)
}
