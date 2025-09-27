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

// GetUser 获取当前用户
func (s *sUser) GetUser(ctx context.Context) (*do.Users, error) {
	userId := s.GetUserId(ctx)
	if userId == "" {
		return nil, nil
	}
	var user *do.Users
	err := dao.Users.Ctx(ctx).Where("uuid = ?", userId).Scan(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserId 从JWT token获取用户ID
func (s *sUser) GetUserId(ctx context.Context) string {
	jwtKey := consts_user.USER_JWT_KEY
	token := g.RequestFromCtx(ctx).GetQuery("token").String()
	if token == "" {
		return ""
	}

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
