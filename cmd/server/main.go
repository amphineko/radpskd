package main

import (
	"flag"
	"time"

	"github.com/amphineko/radpskd/lib/radius"
	"github.com/amphineko/radpskd/lib/redis"
)

func main() {
	var radiusListenAddr string
	radiusListenAddrDefault := "0.0.0.0:1812"
	radiusListenAddrUsage := "Address to listen for RADIUS requests"
	flag.StringVar(&radiusListenAddr, "radius-listen-addr", radiusListenAddrDefault, radiusListenAddrUsage)
	flag.StringVar(&radiusListenAddr, "l", radiusListenAddrDefault, radiusListenAddrUsage)

	var radiusSecret string
	radiusSecretDefault := ""
	radiusSecretUsage := "Shared secret for RADIUS requests"
	flag.StringVar(&radiusSecret, "radius-secret", radiusSecretDefault, radiusSecretUsage)
	flag.StringVar(&radiusSecret, "k", radiusSecretDefault, radiusSecretUsage)

	var redisAddr string
	redisAddrDefault := "127.0.0.1:6379"
	redisAddrUsage := "Address of Redis server"
	flag.StringVar(&redisAddr, "redis-addr", redisAddrDefault, redisAddrUsage)
	flag.StringVar(&redisAddr, "r", redisAddrDefault, redisAddrUsage)

	var redisTTL time.Duration
	redisTTLDefault := 24 * time.Hour
	redisTTLUsage := "TTL for PSK entries"
	flag.DurationVar(&redisTTL, "redis-ttl", redisTTLDefault, redisTTLUsage)
	flag.DurationVar(&redisTTL, "t", redisTTLDefault, redisTTLUsage)

	flag.Parse()

	if radiusSecret == "" {
		panic("Shared secret must be provided")
	}

	redis := redis.New(redisAddr, redisTTL)
	radius := radius.New(radius.ServerConfig{
		Addr:   radiusListenAddr,
		Redis:  redis,
		Secret: radiusSecret,
	})
	if err := radius.ListenAndServe(); err != nil {
		panic(err)
	}
}
