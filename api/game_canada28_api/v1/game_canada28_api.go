package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// EmptyRes 自定义空响应结构体
type EmptyRes struct{}

type GameCommonRes struct{}

// Canada28BetReq 游戏用户acc接口api
type Canada28BetReq struct {
	g.Meta    `path:"/canada28/bet" method:"post" tags:"GameCanada28Bet" summary:"Canada28Bet api"`
	BetType   string  `v:"required"  dc:"bet type"`
	Username  string  `v:"required"  dc:"username"`
	Avatar    string  ``
	UserId    string  `v:"required"  dc:"user id"`
	BetAmount float64 `v:"required"  dc:"bet amount"`
}
