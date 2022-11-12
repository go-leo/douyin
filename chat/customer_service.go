package chat

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"
)

type CustomerServiceData struct {
	Url string `json:"url"` // 官方客服链接
}

type CustomerServiceResp struct {
	Data   *CustomerServiceData `json:"data"`
	ErrNo  int                  `json:"err_no"`
	ErrMsg string               `json:"err_msg"`
	LogId  string               `json:"log_id"`
	AppID  string
}

// CustomerServiceReq 客服请求参数
type CustomerServiceReq struct {
	AccessToken string `json:"access_token"` // 必选,
	OpenID      string `json:"open_id"`      // 必选, 用户在当前小程序的 ID，使用 code2session 接口返回的 openid
	Type        string `json:"type"`         // 必选, 来源，抖音传 1128，抖音极速版传 2329
	Scene       string `json:"scene"`        // 必选, 场景值，固定值 1
	OrderID     string `json:"order_id"`     // 可选, 订单号
	ImType      string `json:"im_type"`      // 必选, im类型, "group_buy"：酒旅、美食、其他本地服务
}

func (sdk *SDK) CustomerService(ctx context.Context, req *CustomerServiceReq) (*CustomerServiceResp, error) {
	var resp CustomerServiceResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(sdk.getCustomerServiceURL()).
		Query("appid", sdk.AppID).
		Query("openid", req.OpenID).
		Query("type", req.Type).
		Query("scene", req.Scene).
		Query("order_id", req.OrderID).
		Query("im_type", req.ImType).
		Header("Access-Token", req.AccessToken).
		Execute(ctx, sdk.HttpCli).JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrNo != 0 {
		err := fmt.Errorf("sns.JsCode2Session error : errcode=%v , errmsg=%v", resp.ErrNo, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = sdk.AppID
	return &resp, nil
}

func (sdk *SDK) getCustomerServiceURL() string {
	return "https://developer.toutiao.com/api/apps/chat/customer_service_url"
}
