package jobs

import (
	"errors"

	"github.com/xxzhwl/gaia"
	"github.com/xxzhwl/gaia/components/jobs"
	"github.com/xxzhwl/gaia/framework/messageImpl"

	"time"
)

func init() {
	jobs.RegisterCronService("firstService", &FirstService{})
	jobs.RegisterCronService("secondService", &SecondService{})
}

type FirstService struct {
}

func (f *FirstService) Start() {
	gaia.Info("CronJobStartFirstService")
}

func (f *FirstService) Stop(arg any) {
	gaia.InfoF("CronJobStopFirstService %v", arg)
}

func (f *FirstService) ReportErr() error {
	return errors.New("CronJobReportErrFirstService")
}

type SecondService struct {
}

func (s *SecondService) End() {
	gaia.Info("CronJobEndSecondService")
	robot := messageImpl.NewFeiShuRobot()
	robot.SendText("我发个消息不过分吧")
	time.Sleep(5 * time.Second)
}

func (s *SecondService) Back() {
	for {
		select {
		case <-time.After(time.Second * 5):
			gaia.Info("这是应该后台进程任务")
		}
	}
}
