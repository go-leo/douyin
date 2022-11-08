package douyin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-leo/backoffx"
	"github.com/go-leo/mapx"
	"github.com/go-leo/netx/httpx"
	"github.com/go-redsync/redsync/v4"
)

type TokenData struct {
	AccessToken string `json:"access_token"` // string 获取到的凭证
	ExpiresIn   int64  `json:"expires_in"`   // number	凭证有效时间，单位：秒
}

type TokenResp struct {
	ErrNo   int64      `json:"err_no"`
	ErrTips string     `json:"err_tips"`
	Data    *TokenData `json:"data"`
}

func (sdk *SDK) GetToken(ctx context.Context) (*TokenResp, error) {
	// 从redis获取token信息，如果获取到了，判断expires_at是否在当前时间之后，如果是之后，则直接返回响应
	// 其他情况需要获取分布式锁，获取到锁的话，就请求微信服务，获取新的token，并更新token，expires_at 设置成的 now+expires_in*(2/3)
	// 没获取到锁的，返回响应（新老有5分钟的过度期）
	result, err := sdk.RedisCli.HGetAll(ctx, sdk.TokenKey).Result()
	if err != nil {
		return nil, err
	}
	if mapx.IsNotEmpty(result) {
		// 从缓存中获取到token信息
		expiresAtTimestamp, _ := result["expires_at"]
		timestamp, _ := strconv.ParseInt(expiresAtTimestamp, 10, 64)
		expiresAt := time.Unix(timestamp, 0)
		if expiresAt.After(time.Now()) {
			// 没有过期，直接返回token信息
			return sdk.DecodeTokenResp(result), nil
		} else {
			// 过期了，获取锁
			mutex := sdk.RedisSync.NewMutex(sdk.TokenLockerKey)
			err := mutex.Lock()
			if err != nil {
				// 获取锁失败，直接返回token信息
				return sdk.DecodeTokenResp(result), nil
			}
			defer func(mutex *redsync.Mutex) {
				_, _ = mutex.Unlock()
			}(mutex)

			// 获取锁成功, 调微信的接口
			tokenResp, err := sdk.CallToken(ctx)
			if err != nil {
				sdk.Logger.Errorf("failed to get access token from wechat, %v", err)
				return sdk.DecodeTokenResp(result), nil
			}

			// 保存到redis
			if err := sdk.SaveTokenRespToRedis(ctx, tokenResp); err != nil {
				sdk.Logger.Errorf("failed to save access token to redis, %v", err)
				return tokenResp, nil
			}
			return tokenResp, nil
		}
	}
	// 从缓存中没有获取到，第一次请求微信获取token
	// 获取锁
	mutex := sdk.RedisSync.NewMutex(
		sdk.TokenLockerKey,
		redsync.WithTries(3),
		redsync.WithRetryDelayFunc(func(tries int) time.Duration {
			return backoffx.Linear(50*time.Millisecond)(ctx, uint(tries))
		}),
	)
	if err := mutex.Lock(); err != nil {
		// 获取锁失败，在从redis中获取一次
		result, err := sdk.RedisCli.HGetAll(ctx, sdk.TokenKey).Result()
		if err != nil {
			return nil, err
		}
		if mapx.IsEmpty(result) {
			return nil, errors.New("failed to get access token")
		}
		return sdk.DecodeTokenResp(result), nil
	}
	defer func(mutex *redsync.Mutex) {
		_, _ = mutex.Unlock()
	}(mutex)

	// 获取锁成功, 调微信的接口
	tokenResp, err := sdk.CallToken(ctx)
	if err != nil {
		sdk.Logger.Errorf("failed to get access token from wechat, %v", err)
		return sdk.DecodeTokenResp(result), nil
	}

	// 保存到redis
	if err := sdk.SaveTokenRespToRedis(ctx, tokenResp); err != nil {
		sdk.Logger.Errorf("failed to save access token to redis, %v", err)
		return tokenResp, nil
	}
	return tokenResp, nil
}

func (sdk *SDK) SaveTokenRespToRedis(ctx context.Context, tokenResp *TokenResp) error {
	data, _ := json.Marshal(tokenResp)
	expiresIn := time.Duration(tokenResp.Data.ExpiresIn) * time.Second
	expiresAt := time.Now().Add(expiresIn * 2 / 3)
	_, err := sdk.RedisCli.HMSet(ctx, sdk.TokenKey, "resp", string(data), "expires_at", expiresAt.Unix()).Result()
	if err != nil {
		return err
	}
	_, err = sdk.RedisCli.Expire(ctx, sdk.TokenKey, expiresIn).Result()
	if err != nil {
		return err
	}
	return nil
}

func (sdk *SDK) DecodeTokenResp(result map[string]string) *TokenResp {
	resp, _ := result["resp"]
	tokenResp := &TokenResp{}
	_ = json.Unmarshal([]byte(resp), tokenResp)
	return tokenResp
}

func (sdk *SDK) CallToken(ctx context.Context) (*TokenResp, error) {
	var resp TokenResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(sdk.getURLToken()).
		Query("appid", sdk.AppID).
		Query("secret", sdk.Secret).
		Query("grant_type", "client_credential").
		Execute(ctx, sdk.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrNo != 0 {
		err = fmt.Errorf("token error : errcode=%v , errmsg=%v", resp.ErrNo, resp.ErrTips)
		return nil, err
	}
	return &resp, nil
}

func (sdk *SDK) getURLToken() string {
	if sdk.IsSandBox {
		return "https://open-sandbox.douyin.com/api/apps/v2/token"
	}
	return "https://developer.toutiao.com/api/apps/v2/token"
}
