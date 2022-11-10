package sns

import (
	"context"
	"fmt"
	"github.com/go-leo/douyin/common"

	"github.com/go-leo/netx/httpx"
)

// JsCode2SessionResp 登录凭证校验的返回结果
type JsCode2SessionResp struct {
	common.BaseResp
	OpenID          string `json:"openid"`           // 用户在当前小程序的 ID，如果请求时有 code 参数才会返回
	SessionKey      string `json:"session_key"`      // 会话密钥，如果请求时有 code 参数才会返回
	AnonymousOpenid string `json:"anonymous_openid"` // 匿名用户在当前小程序的 ID，如果请求时有 anonymous_code 参数才会返回
	UnionID         string `json:"unionid"`          // 用户在小程序平台的唯一标识符，请求时有 code 参数才会返回。如果开发者拥有多个小程序，可通过 unionid 来区分用户的唯一性。
}

func (auth *SDK) JsCode2Session(ctx context.Context, code string) (*JsCode2SessionResp, error) {
	var resp JsCode2SessionResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(URLJsCode2Session).
		Query("appid", auth.AppID).
		Query("secret", auth.Secret).
		Query("code", code).
		Query("anonymous_code", ""). //TODO login 接口返回的匿名登录凭证,和 code至少有一个,这里填写什么
		Execute(ctx, auth.HttpCli).JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err := fmt.Errorf("sns.JsCode2Session error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = auth.AppID
	return &resp, nil
}
