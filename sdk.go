package douyin

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
)

type SDK struct {
	HttpCli        *http.Client
	AppID          string
	Secret         string
	RedisCli       redis.UniversalClient
	RedisSync      *redsync.Redsync
	TokenKey       string
	TokenLockerKey string
	Logger         Logger
	IsSandBox      bool
}
