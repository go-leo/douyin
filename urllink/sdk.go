package urllink

import (
	"net/http"
)

type SDK struct {
	HttpCli   *http.Client
	AppID     string
	IsSandBox bool
}
