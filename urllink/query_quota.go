package urllink

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"
)

type UrlLinkQuota struct {
	// UrlLinkUsed url_link 已生成次数
	UrlLinkUsed string `json:"app_name"`
	// UrlLinkLimit url_link 生成次数上限
	UrlLinkLimit string `json:"ma_app_id"`
}

type QueryQuotaResp struct {
	ErrNo        int           `json:"err_no"`
	ErrTips      string        `json:"err_tips"`
	UrlLinkQuota *UrlLinkQuota `json:"url_link_quota"`
}

type QueryQuotaReq struct {
	AccessToken string `json:"access_token"` // 必选, 小程序 access_token
	MaAppId     string `json:"ma_app_id"`    // 必选, 小程序的 id
}

// QueryQuota 该接口用于查询当前小程序配额。
func (sdk *SDK) QueryQuota(ctx context.Context, req *QueryQuotaReq) (*QueryQuotaResp, error) {
	var resp QueryQuotaResp
	err := httpx.NewRequestBuilder().
		Post().
		URLString(sdk.getQueryQuotaURL()).
		JSONBody(req).
		Execute(ctx, sdk.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrNo != 0 {
		err = fmt.Errorf("QueryInfo error : errcode=%d , errmsg=%s", resp.ErrNo, resp.ErrTips)
		return nil, err
	}
	return &resp, nil
}

func (sdk *SDK) getQueryQuotaURL() string {
	return "https://developer.toutiao.com/api/apps/url_link/query_quota"
}
