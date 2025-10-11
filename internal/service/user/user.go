// internal/service/user.go
package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"demo/internal/dao"
	"demo/internal/model/do"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	consts_user "demo/internal/consts/user"

	"github.com/gogf/gf/v2/frame/g"
)

type sUser struct{}

var User = sUser{}

// JWTPayload JWT载荷结构
type JWTPayload struct {
	UserID string `json:"user_id"`
	Exp    int64  `json:"exp"`
}

// GetUser 获取当前用户（支持普通用户和admin）
func (s *sUser) GetUser(ctx context.Context) (*do.Users, error) {
	userId := s.GetUserId(ctx)
	if userId == "" {
		return nil, nil
	}

	// 先尝试从game_users表查询
	var user *do.Users
	err := dao.Users.Ctx(ctx).Where("uuid = ?", userId).Scan(&user)
	if err != nil {
		return nil, err
	}

	// 如果在game_users表中找到了，直接返回
	if user != nil {
		return user, nil
	}

	// 如果没找到，尝试从game_admin表查询（admin用户）
	var admin *do.Admin
	err = dao.Admin.Ctx(ctx).Where("uuid = ?", userId).Scan(&admin)
	if err != nil {
		return nil, err
	}

	// 如果找到admin，转换为Users结构返回
	if admin != nil {
		user = &do.Users{
			Uuid:     admin.Uuid,
			Nickname: admin.Nickname,
			Username: admin.Username,
			Avatar:   admin.Avatar,
			Status:   admin.Status,
		}
		return user, nil
	}

	return nil, nil
}

// GetUserId 从JWT token或UUID获取用户ID
func (s *sUser) GetUserId(ctx context.Context) string {
	token := g.RequestFromCtx(ctx).GetQuery("token").String()
	if token == "" {
		return ""
	}

	// 检查是否为UUID格式（admin token）
	// UUID格式: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (36字符，包含4个连字符)
	if len(token) == 36 && strings.Count(token, "-") == 4 {
		// 直接返回UUID作为用户ID（admin）
		return token
	}

	// JWT token处理（普通用户）
	jwtKey := consts_user.USER_JWT_KEY

	// 分割token
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return ""
	}

	payload := parts[0]
	signature := parts[1]

	// 验证签名
	h := hmac.New(sha256.New, []byte(jwtKey))
	h.Write([]byte(payload))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != signature {
		return ""
	}

	// 解码payload
	payloadBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return ""
	}

	var data JWTPayload
	err = json.Unmarshal(payloadBytes, &data)
	if err != nil {
		return ""
	}

	// 检查必要字段
	if data.UserID == "" || data.Exp == 0 {
		return ""
	}

	// 检查过期时间
	if data.Exp < time.Now().Unix() {
		return ""
	}

	return data.UserID
}
