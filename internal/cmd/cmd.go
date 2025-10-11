package cmd

import (
	"context"
	"demo/boot"
	"demo/internal/controller/customer_service_ws"
	"demo/internal/controller/game_bingo28_api"
	"demo/internal/controller/game_bingo28_ws"
	"demo/internal/controller/game_canada28_api"
	"demo/internal/controller/game_canada28_ws"
	"demo/internal/controller/game_keno_ws"
	"demo/internal/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gtime"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			boot.InitLogger(ctx)
			// 初始化OSS
			boot.InitializeOSS()
			// 初始化RabbitMQ
			// boot.InitRabbitMQ(ctx)

			s := g.Server()
			gtime.SetTimeZone("UTC")

			// 初始化定时任务
			RegisterCron(ctx)

			// 注册全局中间件
			s.Use(middleware.TraceMiddleware)

			s.AddStaticPath("/uploads", "/uploads")

			// 全局跨域中间件
			s.Use(ghttp.MiddlewareCORS)

			// 游戏api路由
			s.Group("/v1/game_api", func(group *ghttp.RouterGroup) {
				// 使用自定义的响应处理中间件，替代默认的响应处理中间件
				group.Middleware(middleware.ResponseHandler)
				group.Bind(
					game_canada28_api.NewV1(),
					game_bingo28_api.NewV1(),
				)
			})

			// WebSocket路由
			s.Group("/game_canada28_ws", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.NewUserAuth().UserAuth)
				group.GET("/connect", game_canada28_ws.NewWsController().Connect)
			})
			// WebSocket路由
			s.Group("/game_bingo28_ws", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.NewUserAuth().UserAuth)
				group.GET("/connect", game_bingo28_ws.NewWsController().Connect)
			})
			// WebSocket路由
			s.Group("/game_keno_ws", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.NewUserAuth().UserAuth)
				group.GET("/connect", game_keno_ws.NewWsController().Connect)
			})
			// Customer Service WebSocket路由
			s.Group("/customer_service_ws", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.NewUserAuth().UserAuth)
				group.GET("/connect", customer_service_ws.NewWsController().Connect)
			})
			s.Run()
			return nil
		},
	}
)
