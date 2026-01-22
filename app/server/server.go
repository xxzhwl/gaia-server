// Package server 包注释
// @author wanlizhan
// @created 2025-03-17
package server

import (
	"gaia-server/app/server/asynctask"
	_ "gaia-server/app/server/jobs"
	"gaia-server/app/server/repo/dao"
	"gaia-server/app/server/router"

	"github.com/xxzhwl/gaia"
	"github.com/xxzhwl/gaia/framework"
	"github.com/xxzhwl/gaia/framework/server"
)

func init() {
	asynctask.RegisterTasks()
}

func RunServer() {
	framework.Init()
	if err := defaultDb(); err != nil {
		gaia.PanicLog(err.Error())
	}
	s := server.DefaultApp()

	router.RegisterRouters(s)
	s.Run()
}

func defaultDb() error {
	db, err := gaia.NewFrameworkMysql()
	if err != nil {
		return err
	}
	dao.SetDefault(db.GetGormDb())
	return nil
}
