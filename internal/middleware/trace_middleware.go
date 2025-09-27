package middleware

import (
	"demo/internal/consts/ctx_var"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// TraceMiddleware 为请求添加追踪功能的中间件
func TraceMiddleware(r *ghttp.Request) {
	// 请求开始时间
	startTime := time.Now()
	ctx := r.GetCtx()

	// 记录请求信息
	logMap := g.Map{
		"action":  "request",
		"method":  r.Method,
		"path":    r.URL.Path,
		"ip":      r.GetClientIp(),
		"headers": r.Header,
	}
	// 安全地获取参数，避免ParseForm失败导致的错误
	params := r.GetMap()
	if len(params) > 0 {
		logMap["params"] = params
		r.SetCtxVar(ctx_var.REQUEST_CONTENT, params)
	}

	g.Log().Info(ctx, logMap)

	// 修改响应处理方式，包装响应内容
	r.Middleware.Next()

	// 请求结束，记录响应信息
	duration := time.Since(startTime)
	responseData := r.Response.BufferString()

	// 尝试解析JSON响应并添加TraceID
	if r.Response.Header().Get("Content-Type") == "application/json" && responseData != "" {
		var responseMap map[string]interface{}
		if err := json.Unmarshal([]byte(responseData), &responseMap); err == nil {
			// 只有在JSON解析成功的情况下才修改响应
			if newResponseData, err := json.Marshal(responseMap); err == nil {
				// 重置响应内容
				r.Response.ClearBuffer()
				r.Response.Write(newResponseData)
				// 更新响应数据变量
				responseData = r.Response.BufferString()
			}
		}
	}

	g.Log().Info(ctx, g.Map{
		"action":   "response",
		"method":   r.Method,
		"path":     r.URL.Path,
		"status":   r.Response.Status,
		"duration": duration.String(),
		"response": responseData,
	})
}
