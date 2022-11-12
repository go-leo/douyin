package douyin

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"
)

type JsCode2SessionData struct {
	SessionKey      string `json:"session_key"`
	Openid          string `json:"openid"`
	AnonymousOpenid string `json:"anonymous_openid"`
	Unionid         string `json:"unionid"`
}

// JsCode2SessionResp 登录凭证校验的返回结果
type JsCode2SessionResp struct {
	ErrNo   int                 `json:"err_no"`
	ErrTips string              `json:"err_tips"`
	Data    *JsCode2SessionData `json:"data"`
	AppID   string
}

func (sdk *SDK) JsCode2Session(ctx context.Context, code string, anonymousCode string) (*JsCode2SessionResp, error) {
	var resp JsCode2SessionResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(sdk.getJsCode2SessionUrl()).
		Query("appid", sdk.AppID).
		Query("secret", sdk.Secret).
		Query("code", code).
		Query("anonymous_code", anonymousCode).
		Execute(ctx, sdk.HttpCli).JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrNo != 0 {
		err := fmt.Errorf("JsCode2Session error : errcode=%v , errmsg=%v", resp.ErrNo, resp.ErrTips)
		return nil, err
	}
	resp.AppID = sdk.AppID
	return &resp, nil
}

func (sdk *SDK) getJsCode2SessionUrl() string {
	if sdk.IsSandBox {
		return "https://open-sandbox.douyin.com/api/apps/v2/jscode2session"
	}
	return "https://developer.toutiao.com/api/apps/v2/jscode2session"
}
