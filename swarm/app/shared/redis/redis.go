package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisInterface interface {
	Close()

	Keys(ctx context.Context, pattern string) []string
	Get(ctx context.Context, key string) string
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration)
	Del(ctx context.Context, keys ...string)

	Publish(ctx context.Context, channel string, message interface{})
	PSubscribe(ctx context.Context, channels ...string) *redis.PubSub
}

type Redis struct {
	client redis.UniversalClient
}

func New(ctx context.Context, appID string) RedisInterface {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Panic().Err(err).Msg("redis.New")
	}

	client := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		ClientName: appID,
		DB:         db,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Panic().Err(err).Msg("redis.New")
	}

	if err := redisotel.InstrumentTracing(client); err != nil {
		log.Panic().Err(err).Msg("redis.New")
	}

	if err := redisotel.InstrumentMetrics(client); err != nil {
		log.Panic().Err(err).Msg("redis.New")
	}

	return &Redis{client}
}

func (r *Redis) Close() {
	if err := r.client.Close(); err != nil {
		log.Error().Err(err).Msg("redis.Close")
	}
}

func (r *Redis) Keys(ctx context.Context, pattern string) []string {
	v, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		log.Error().Err(err).Msg("redis.Keys")
	}
	return v
}

func (r *Redis) Get(ctx context.Context, key string) string {
	v, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Msg("redis.Get")
	}
	return v
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	if err := r.client.Set(ctx, key, value, expiration).Err(); err != nil {
		log.Error().Err(err).Msg("redis.Set")
	}
}

func (r *Redis) Del(ctx context.Context, keys ...string) {
	k := []string{}
	for _, key := range keys {
		if strings.Contains(key, "*") {
			v, _ := r.client.Keys(ctx, key).Result()
			k = append(k, v...)
		} else {
			k = append(k, key)
		}
	}
	if err := r.client.Del(ctx, k...).Err(); err != nil {
		log.Error().Err(err).Msg("redis.Del")
	}
}

func (r *Redis) Publish(ctx context.Context, channel string, message interface{}) {
	if err := r.client.Publish(ctx, channel, message).Err(); err != nil {
		log.Error().Err(err).Msg("redis.Publish")
	}
}

func (r *Redis) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return r.client.PSubscribe(ctx, channels...)
}
