package testutil

import (
	"testing"

	"example.com/shared/redis"
)

func TestRedis(t *testing.T) {
	var v redis.RedisInterface
	v = &Redis{}
	_ = v
}
