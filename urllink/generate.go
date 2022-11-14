package urllink

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"
)

type GenerateResp struct {
	ErrNo   int64  `json:"err_no"`
	ErrTips string `json:"err_tips"`
	UrlLink string `json:"url_link"`
}

type GenerateReq struct {
	AccessToken string `json:"access_token"` // 必选, 小程序 access_token
	MaAppId     string `json:"ma_app_id"`    // 必选, 小程序的 id
	AppName     string `json:"app_name"`     // 必选, 宿主名称，可选 douyin，douyinlite
	Path        string `json:"path"`         // 可选, 通过URL Link进入的小程序页面路径，必须是已经发布的小程序存在的页面，不可携带 query。path 为空时会跳转小程序主页。
	Query       string `json:"query"`        // 可选, 通过URL Link进入小程序时的 query（json形式），若无请填{}。最大1024个字符，只支持数字，大小写英文以及部分特殊字符：`{}!#$&'()*+,/:;=?@-._~%``。
	ExpireTime  string `json:"expire_time"`  // 必选, 到期失效的URL Link的失效时间。为 Unix 时间戳，实际失效时间为距离当前时间小时数，向上取整。最长间隔天数为180天。
}

// Generate 该接口用于生成能够直接跳转到端内小程序的 url link。
func (sdk *SDK) Generate(ctx context.Context, req *GenerateReq) (*GenerateResp, error) {
	var resp GenerateResp
	err := httpx.NewRequestBuilder().
		Post().
		URLString(sdk.getGenerateURL()).
		JSONBody(req).
		Execute(ctx, sdk.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrNo != 0 {
		err = fmt.Errorf("Generate error : errcode=%d , errmsg=%s", resp.ErrNo, resp.ErrTips)
		return nil, err
	}
	return &resp, nil
}

func (sdk *SDK) getGenerateURL() string {
	return "https://developer.toutiao.com/api/apps/url_link/generate"
}
