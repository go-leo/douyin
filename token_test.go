package douyin

import (
	"context"
	"testing"

	"github.com/go-leo/netx/httpx"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/stretchr/testify/assert"
)

func TestSDK_Token(t *testing.T) {

	redisCli := ""
	password := ""
	appID := ""
	secret := ""

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{redisCli},
		Password: password,
	})
	tokenKey := "access:token:dy:" + appID
	tokenLockerKey := tokenKey + ":locker"
	sdk := &SDK{
		HttpCli:        httpx.PooledClient(),
		AppID:          appID,
		Secret:         secret,
		RedisCli:       client,
		RedisSync:      redsync.New(goredis.NewPool(client)),
		TokenKey:       tokenKey,
		TokenLockerKey: tokenLockerKey,
		IsSandBox:      false,
		Logger:         DefaultLogger{},
	}
	tokenResp, err := sdk.Token(context.Background())
	assert.NoError(t, err)
	t.Log(tokenResp)
}
