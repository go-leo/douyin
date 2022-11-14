package urllink

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"
)

type QueryInfoInfo struct {
	// AppName 宿主名称，douyin，douyinlite
	AppName string `json:"app_name"`
	// MaAppId 小程序ID
	MaAppId string `json:"ma_app_id"`
	// Path 小程序页面路径。
	Path string `json:"path"`
	// Query 小程序页面query。
	Query string `json:"query"`
	// CreateTime 创建时间，为 Unix 时间戳
	CreateTime int `json:"create_time"`
	// ExpireTime 到期失效时间，为 Unix 时间戳
	ExpireTime int `json:"expire_time"`
}

type QueryInfoResp struct {
	ErrNo       int            `json:"err_no"`
	ErrTips     string         `json:"err_tips"`
	UrlLinkInfo *QueryInfoInfo `json:"url_link_info"`
}

type QueryInfoReq struct {
	AccessToken string `json:"access_token"` // 必选, 小程序 access_token
	MaAppId     string `json:"ma_app_id"`    // 必选, 小程序的 id
	UrlLink     string `json:"url_link"`     // 必选, 生成的url_link
}

// QueryInfo 该接口用于查询已经生成的 link 的信息。
func (sdk *SDK) QueryInfo(ctx context.Context, req *QueryInfoReq) (*QueryInfoResp, error) {
	var resp QueryInfoResp
	err := httpx.NewRequestBuilder().
		Post().
		URLString(sdk.getQueryInfoURL()).
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

func (sdk *SDK) getQueryInfoURL() string {
	return "https://developer.toutiao.com/api/apps/url_link/query_info"
}
