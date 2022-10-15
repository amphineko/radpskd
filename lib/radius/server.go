package radius

import (
	"github.com/amphineko/radpskd/lib/log"
	"github.com/amphineko/radpskd/lib/redis"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/vendors/aruba"
)

type ServerConfig struct {
	Addr   string
	Redis  *redis.RedisContext
	Secret string
}

type RadiusContext struct {
	server radius.PacketServer
	redis  *redis.RedisContext
}

func New(config ServerConfig) *RadiusContext {
	ctx := &RadiusContext{
		redis: config.Redis,
	}
	ctx.server = radius.PacketServer{
		Addr:         config.Addr,
		Handler:      radius.HandlerFunc(ctx.HandleRadiusRequest),
		SecretSource: radius.StaticSecretSource([]byte(config.Secret)),
	}
	return ctx
}

func (ctx *RadiusContext) HandleRadiusRequest(w radius.ResponseWriter, r *radius.Request) {
	hwaddr := rfc2865.CallingStationID_GetString(r.Packet)
	log.Debug.Printf("[radius.HandleRadiusRequest] processing hwaddr: %s", hwaddr)

	psk, err := ctx.redis.GetHwaddrPsk(hwaddr)
	if err != nil {
		log.Error.Printf("[radius.HandleRadiusRequest] failed to get psk for hwaddr %s: %s", hwaddr, err)
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	resp := r.Response(radius.CodeAccessAccept)
	if err := aruba.ArubaMPSKPassphrase_SetString(resp, psk); err != nil {
		log.Error.Printf("[radius.HandleRadiusRequest] failed to write ArubaMPSKPassphrase: %s", err)
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	log.Info.Printf("[radius.HandleRadiusRequest] accepting hwaddr %s with psk: %s", hwaddr, psk)
	w.Write(resp)
}

func (ctx *RadiusContext) ListenAndServe() error {
	return ctx.server.ListenAndServe()
}
