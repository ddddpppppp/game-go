package game_canada28_api

import (
	"context"
	v1 "demo/api/game_api/v1"
	"demo/internal/service/game_canada28_ws"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func (c *ControllerV1) Canada28Bet(ctx context.Context, req *v1.Canada28BetReq) (res *v1.GameCommonRes, err error) {

	message := fmt.Sprintf("bet %s, amount $%.2f", req.BetType, req.BetAmount)
	err = game_canada28_ws.BroadcastToAllUsers(ctx, g.Map{
		"action": "new_message",
		"data": g.Map{
			"id":         "",
			"nickname":   req.Username,
			"avatar":     req.Avatar,
			"user_id":    req.UserId,
			"type":       "text",
			"message":    message,
			"created_at": gtime.Now(),
		},
	})

	if err != nil {
		return &v1.GameCommonRes{}, err
	}

	message = fmt.Sprintf("@%s bet %s successfully, amount: $%.2f", req.Username, req.BetType, req.BetAmount)
	err = game_canada28_ws.BroadcastToAllUsers(ctx, g.Map{
		"action": "new_message",
		"data": g.Map{
			"id":         "",
			"nickname":   "bot",
			"avatar":     "",
			"user_id":    "bot",
			"type":       "text",
			"message":    message,
			"created_at": gtime.Now(),
		},
	})

	if err != nil {
		return &v1.GameCommonRes{}, err
	}
	return &v1.GameCommonRes{}, nil
}
