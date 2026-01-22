package asynctask

import (
	"errors"
	"time"

	"github.com/xxzhwl/gaia"
	"github.com/xxzhwl/gaia/framework/messageImpl"
)

func RegisterTasks() {
	gaia.RegisterProxy("asynctasks", "firstService", &FirstService{})
	gaia.RegisterProxy("asynctasks", "secondService", &SecondService{})
}

type FirstService struct {
}

func (f *FirstService) Start() {
	time.Sleep(2 * time.Second)
	gaia.Info("AsyncTaskStartFirstService")
	messageImpl.NewFeiShuRobot().SendRichText(
		"我试一下",
		[]messageImpl.Content{
			{
				Tag:  "text",
				Text: "我试一下",
			},
		},
	)
}

func (f *FirstService) Stop(arg any) {
	gaia.InfoF("AsyncTaskStopFirstService %v", arg)
}

func (f *FirstService) ReportErr() error {
	return errors.New("AsyncTaskReportErrFirstService")
}

type SecondService struct {
}

func (f *SecondService) End() {
	gaia.Info("AsyncTaskEndSecondService")
}
