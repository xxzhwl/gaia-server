package api

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/framework/server"
)

type DemoCtrl struct{}

func NewDemoCtrl() *DemoCtrl {
	return &DemoCtrl{}
}

type DemoRequest struct {
	Name string `json:"name" require:"1"`
}

type DemoResponse struct {
	Name string `json:"name"`
}

func (d *DemoCtrl) Demo() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (DemoResponse, error) {
		req := DemoRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return DemoResponse{}, err
		}
		return DemoResponse{Name: req.Name}, nil
	})
}
