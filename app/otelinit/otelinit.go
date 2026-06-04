// Package otelinit 在框架初始化之前设置 OTel 全局错误处理器，
// 避免 OTel 组件在初始化过程中向 stderr 打印 "<nil>" 日志。
package otelinit

import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
)

func init() {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		if err == nil {
			return
		}
		// 框架日志尚未就绪，真实错误暂时输出到 stderr
		fmt.Fprintf(os.Stderr, "[OTel] %v\n", err)
	}))
}
