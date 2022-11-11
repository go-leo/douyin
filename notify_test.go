package douyin

import (
	"context"
	"github.com/go-leo/douyin/common"
	"testing"

	"github.com/go-leo/netx/httpx"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func TestNotify(t *testing.T) {
	appID := ""
	secret := ""
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{""},
		Password: "",
	})
	sdk := &SDK{
		HttpCli:        httpx.PooledClient(),
		AppID:          appID,
		Secret:         secret,
		RedisCli:       client,
		RedisSync:      redsync.New(goredis.NewPool(client)),
		TokenKey:       "douying:access:token:" + appID,
		TokenLockerKey: "douying:access:token:locker" + appID,
		Logger:         common.DefaultLogger{},
		IsSandBox:      true,
	}
	token, err := sdk.GetToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token.Data.AccessToken)
}
