package storage

import (
	"context"
	"time"

	"example.com/shared/redis"
)

type Storage struct {
	redis  redis.RedisInterface
	prefix string
}

func New(redis redis.RedisInterface, prefix string) *Storage {
	return &Storage{redis, prefix}
}

func (s *Storage) Get(key string) ([]byte, error) {
	v := s.redis.Get(context.Background(), s.prefix+key)
	if v == "" {
		return nil, nil
	}
	return []byte(v), nil
}

func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	s.redis.Set(context.Background(), s.prefix+key, string(val), exp)
	return nil
}

func (s *Storage) Delete(key string) error {
	s.redis.Del(context.Background(), s.prefix+key)
	return nil
}

func (s *Storage) Reset() error {
	s.redis.Del(context.Background(), s.prefix+"*")
	return nil
}

func (s *Storage) Close() error {
	// unused
	return nil
}
