package douyin

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"
)

type NotifyResp struct {
	ErrNo   int64  `json:"err_no"`
	ErrTips string `json:"err_tips"`
}

// NotifyReq 订阅消息请求参数
type NotifyReq struct {
	AccessToken string            `json:"access_token"` // 必选, 小程序 access_token
	AppID       string            `json:"app_id"`       // 必选, 小程序的 id
	TplID       string            `json:"tpl_id"`       // 必选, 模板的 id
	OpenID      string            `json:"open_id"`      // 必选, 接收消息目标用户的 open_id
	Data        map[string]string `json:"data"`         // 必选, 模板内容，格式形如 { "key1": "value1", "key2": "value2" }，具体使用方式参考下文请求示例
	Page        string            `json:"page"`         // 可选, 跳转的页面
}

// Notify 发送订阅消息
func (sdk *SDK) Notify(ctx context.Context, accessToken, tplId, openID, page string, data map[string]string) (*NotifyResp, error) {
	req := NotifyReq{
		AccessToken: accessToken,
		AppID:       sdk.AppID,
		TplID:       tplId,
		OpenID:      openID,
		Data:        data,
		Page:        page,
	}
	var resp NotifyResp
	err := httpx.NewRequestBuilder().
		Post().
		URLString(sdk.getNotifyURL()).
		JSONBody(&req).
		Execute(ctx, sdk.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrNo != 0 {
		err = fmt.Errorf("notify error : errcode=%d , errmsg=%s", resp.ErrNo, resp.ErrTips)
		return nil, err
	}
	return &resp, nil
}

func (sdk *SDK) getNotifyURL() string {
	if sdk.IsSandBox {
		return "https://open-sandbox.douyin.com/api/apps/subscribe_notification/developer/v1/notify"
	}
	return "https://developer.toutiao.com/api/apps/subscribe_notification/developer/v1/notify"
}
