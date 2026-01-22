package api

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/framework/httpclient"
	"github.com/xxzhwl/gaia/framework/server"
)

type DemoCtrl struct{}

func NewDemoCtrl() *DemoCtrl {
	return &DemoCtrl{}
}

type DemoRequest struct {
	Name string `json:"name" require:"1"`
}

func (d *DemoCtrl) Demo() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		req := DemoRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, err
		}
		return map[string]any{"Name": req.Name}, nil
	})
}

func (d *DemoCtrl) Req() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		r := httpclient.NewHttpRequest("http://12.24.1.2/api")
		body, _, err := r.WithTitle("DemoSystem").WithTimeOut(5 * time.Second).Do()
		if err != nil {
			return nil, err
		}
		return map[string]any{"Res": string(body)}, nil
	})
}
