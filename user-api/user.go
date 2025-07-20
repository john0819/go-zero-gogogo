package main

import (
	"flag"
	"fmt"

	"zero-demo/user-api/internal/config"
	"zero-demo/user-api/internal/handler"
	"zero-demo/user-api/internal/middleware"
	"zero-demo/user-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// go-zero会自动根据Telemetry配置初始化OpenTelemetry
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册trace中间件
	server.Use(middleware.TraceMiddleware)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	if c.Telemetry.Endpoint != "" {
		fmt.Printf("OpenTelemetry tracing enabled, endpoint: %s\n", c.Telemetry.Endpoint)
	}
	server.Start()
}
