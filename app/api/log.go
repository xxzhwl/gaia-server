package api

import (
	"errors"
	"gaia-server/app/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/errwrap"
	"github.com/xxzhwl/gaia/framework/server"
)

// LogCtrl 日志查询控制器
type LogCtrl struct{}

func NewLogCtrl() *LogCtrl {
	return &LogCtrl{}
}

// SimpleSearchArg ES查询参数
type SimpleSearchArg struct {
	Index  string           `json:"index" require:"1"`
	Sorts  []service.SortKv `json:"sorts"`
	From   int              `json:"from"`
	Size   int              `json:"size"`
	Must   []service.OpArg  `json:"must"`
	Should []service.OpArg  `json:"should"`
	Not    []service.OpArg  `json:"not"`
}

// QueryES 查询ES日志
func (l *LogCtrl) QueryES() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		req := SimpleSearchArg{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		arg.BanLogger("日志查询", true)

		logService := service.GetLogService()
		result, err := logService.QueryES(req.Index, req.From, req.Size, req.Must, req.Should, req.Not, req.Sorts)
		if err != nil {
			return nil, err
		}

		return result, nil
	})
}

// GetIndices 获取ES索引列表
func (l *LogCtrl) GetIndices() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) ([]string, error) {
		logService := service.GetLogService()
		return logService.GetIndices()
	})
}

// GetMapping 获取ES字段映射
func (l *LogCtrl) GetMapping() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		index := arg.GetUrlQuery("index")
		if index == "" {
			return nil, errwrap.Error(400, errors.New("index参数不能为空"))
		}

		logService := service.GetLogService()
		return logService.GetMapping(index)
	})
}

// ExportLogs 导出日志
func (l *LogCtrl) ExportLogs() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) ([]byte, error) {
		req := SimpleSearchArg{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		arg.BanLogger("日志查询", true)
		logService := service.GetLogService()
		return logService.ExportLogs(req.Index, req.Must, req.Should, req.Not)
	})
}
