package sns

import (
	"net/http"
)

var (
	URLJsCode2Session = "https://developer.toutiao.com/api/apps/v2/jscode2session" //TODO 这里是线上地址
)

type SDK struct {
	HttpCli *http.Client
	AppID   string
	Secret  string
}
