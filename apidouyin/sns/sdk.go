package sns

import (
	"net/http"
)

var (
	URLJsCode2Session = "https://developer.toutiao.com/api/apps/v2/jscode2session"
)

type SDK struct {
	HttpCli *http.Client
	AppID   string
	Secret  string
}
