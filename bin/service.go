// Package main 包注释
// @author wanlizhan
// @created 2025-03-17
package main

import (
	"flag"
	"gaia-server/app/server"

	"github.com/xxzhwl/gaia"
	"github.com/xxzhwl/gaia/framework"
	_ "github.com/xxzhwl/gaia/init"
)

func main() {
	framework.Init()
	s := flag.String("Service", "", "调用的服务名称")
	flag.Parse()
	_, err := gaia.CallMethodWithArgs(&Service{}, *s)
	if err != nil {
		panic(err)
	}
}

type Service struct{}

func (s *Service) Server() {
	server.RunServer()
}

func (s *Service) Test() {
}
