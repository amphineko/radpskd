package redis

import (
	"context"
	"time"

	"github.com/amphineko/radpskd/lib/log"
	"github.com/amphineko/radpskd/lib/utils"
	"github.com/go-redis/redis/v9"
)

type RedisContext struct {
	Ctx        context.Context
	DefaultTtl time.Duration
	Redis      *redis.Client
}

func New(addr string, defaultTtl time.Duration) *RedisContext {
	return &RedisContext{
		Ctx:        context.Background(),
		DefaultTtl: defaultTtl,
		Redis:      redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (ctx *RedisContext) GetHwaddrPsk(hwaddr string) (string, error) {
	if err := ctx.Redis.Expire(ctx.Ctx, hwaddr, ctx.DefaultTtl).Err(); err != nil {
		// renew the key
		log.Error.Printf("[redis.GetHwaddrPsk] failed to renew hwaddr %s: %s", hwaddr, err)
		return "", err
	}

	log.Debug.Printf("[redis.GetHwaddrPsk] reading psk for hwaddr %s", hwaddr)
	psk, err := ctx.Redis.Get(ctx.Ctx, hwaddr).Result()
	if err == redis.Nil {
		// hwaddr not found, generate a new psk
		log.Info.Printf("[redis.GetHwaddrPsk] hwaddr %s not found", hwaddr)
		return ctx.AddHwaddrPsk(hwaddr)
	}
	if err != nil {
		log.Error.Printf("[redis.GetHwaddrPsk] failed to read psk for hwaddr %s: %s", hwaddr, err)
		return "", err
	}

	return psk, nil
}

func (ctx *RedisContext) AddHwaddrPsk(hwaddr string) (string, error) {
	newPsk, err := utils.GeneratePsk()
	if err != nil {
		log.Error.Printf("[redis.AddHwaddrPsk] failed to generate psk for hwaddr %s: %s", hwaddr, err)
		return "", err
	}

	oldPsk, err := ctx.Redis.SetArgs(ctx.Ctx, hwaddr, newPsk, redis.SetArgs{
		Get:  true,
		TTL:  ctx.DefaultTtl,
		Mode: "NX",
	}).Result()

	if oldPsk != "" {
		// hwaddr already exists
		log.Info.Printf("[redis.AddHwaddrPsk] hwaddr %s already exists", hwaddr)
		return oldPsk, nil
	}

	if err != redis.Nil && err != nil {
		log.Error.Printf("[redis.AddHwaddrPsk] failed to add psk for hwaddr %s: %s", hwaddr, err)
		return "", err
	}

	log.Info.Printf("[redis.AddHwaddrPsk] added psk for hwaddr %s: %s", hwaddr, newPsk)
	return newPsk, nil
}
