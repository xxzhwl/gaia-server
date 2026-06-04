// Package main 包注释
// @author wanlizhan
// @created 2025-03-17
package main

import (
	"flag"

	_ "gaia-server/app/otelinit"

	"github.com/xxzhwl/gaia"
	_ "github.com/xxzhwl/gaia/init"

	server "gaia-server/app"
)

func init() {
	gaia.BuildContextTrace()
}

func main() {
	s := flag.String("Service", "", "调用的服务名称")
	a := flag.String("Arg", "", "调用服务的参数")
	flag.Parse()
	_, err := gaia.CallMethodWithArgs(&Service{}, *s, *a)
	if err != nil {
		panic(err)
	}
}

type Service struct{}

func (s *Service) Server(port string) {
	server.RunServer(port)
}
