// internal/middleware/user_auth.go
package middleware

import (
	"demo/internal/model"
	service "demo/internal/service/user"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

type UserAuth struct{}

func NewUserAuth() *UserAuth {
	return &UserAuth{}
}

// UserAuth 用户认证中间件
func (m *UserAuth) UserAuth(r *ghttp.Request) {
	ctx := r.Context()
	// 1. 获取用户信息
	user, err := service.User.GetUser(ctx)
	if err != nil {
		m.authError(r, "Get User Error")
		return
	}

	// if user == nil {
	// 	m.authError(r, "Token is invalid")
	// 	return
	// }
	// if gconv.Int(user.Status) != 1 {
	// 	m.authError(r, "User is disabled")
	// 	return
	// }
	if user != nil {
		// 将用户信息存入上下文
		r.SetCtxVar(model.CtxUserKey, user)
	}

	// 继续执行后续逻辑
	r.Middleware.Next()
}

// authError 统一的认证错误处理
func (m *UserAuth) authError(r *ghttp.Request, message string) {
	ctx := r.Context()
	g.Log().Infof(ctx, "User auth error: %s", message)

	r.Response.WriteHeader(http.StatusUnauthorized)
	r.Response.WriteJson(g.Map{
		"status":     -1,
		"statusText": message,
		"timestamp":  gtime.Now().Unix(),
	})

	// 终止后续处理
	r.Exit()
}
