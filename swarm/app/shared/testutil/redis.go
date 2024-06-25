package testutil

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	KeysN       int
	KeysPattern []string
	KeysReturn  [][]string

	GetN      int
	GetKey    []string
	GetReturn []string

	SetN          int
	SetKey        []string
	SetValue      []interface{}
	SetExpiration []time.Duration

	DelN    int
	DelKeys [][]string

	PublishN       int
	PublishChannel []string
	PublishMessage []interface{}

	PSubscribeN        int
	PSubscribeChannels [][]string
}

func (r *Redis) Close() {}

func (r *Redis) Keys(ctx context.Context, pattern string) []string {
	r.KeysPattern = append(r.KeysPattern, pattern)
	var v []string
	if len(r.KeysReturn) != 0 {
		v = r.KeysReturn[r.KeysN]
	}
	r.KeysN += 1
	return v
}

func (r *Redis) Get(ctx context.Context, key string) string {
	r.GetKey = append(r.GetKey, key)
	var v string
	if len(r.GetReturn) != 0 {
		v = r.GetReturn[r.GetN]
	}
	r.GetN += 1
	return v
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	r.SetKey = append(r.SetKey, key)
	r.SetValue = append(r.SetValue, value)
	r.SetExpiration = append(r.SetExpiration, expiration)
	r.SetN += 1
}

func (r *Redis) Del(ctx context.Context, keys ...string) {
	r.DelKeys = append(r.DelKeys, keys)
	r.DelN += 1
}

func (r *Redis) Publish(ctx context.Context, channel string, message interface{}) {
	r.PublishChannel = append(r.PublishChannel, channel)
	r.PublishMessage = append(r.PublishMessage, message)
	r.PublishN += 1
}

func (r *Redis) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	r.PSubscribeChannels = append(r.PSubscribeChannels, channels)
	r.PSubscribeN += 1
	return nil // TODO
}
