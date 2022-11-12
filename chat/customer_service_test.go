package chat

import (
	"context"
	"testing"

	"github.com/go-leo/netx/httpx"
)

func TestSDK_CustomerService(t *testing.T) {
	appID := ""
	secret := ""
	sdk := &SDK{
		HttpCli:   httpx.PooledClient(),
		AppID:     appID,
		IsSandBox: true,
	}
	service, err := sdk.CustomerService(context.Background(), &CustomerServiceReq{
		OpenID:  "",
		Type:    "1128",
		Scene:   "1",
		OrderID: "",
		ImType:  "group_buy",
	})

}
