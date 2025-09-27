package crontab_game_canada28

import (
	"context"
	consts_sync "demo/internal/consts/sync"
	"demo/internal/dao"
	"demo/internal/model/do"
	"fmt"
	"time"

	game_canada28_ws "demo/internal/service/game_canada28_ws"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type LockCanada28Cron struct {
	ctx context.Context
}

func NewLockCanada28Cron(ctx context.Context) *LockCanada28Cron {
	return &LockCanada28Cron{
		ctx: ctx,
	}
}

// 提前30秒发送锁定通知
func (c *LockCanada28Cron) DoLock() error {
	var draw *do.Canada28Draws
	// 在开奖时间前30秒发送锁定通知
	lockTime := gtime.Now().Add(30 * time.Second)
	err := dao.Canada28Draws.Ctx(c.ctx).Where("status in (0,1) and end_at < ?", lockTime).Order("id desc").Scan(&draw)
	if err != nil {
		return err
	}
	if draw == nil {
		return nil
	}
	// 群发通知
	// c.sendGroupMsg(draw)
	return nil
}

func (c *LockCanada28Cron) sendGroupMsg(draw *do.Canada28Draws) {
	var redisKey = fmt.Sprintf(consts_sync.SyncCanada28AlertMsgKey, draw.Id)
	res, err := g.Redis().Get(c.ctx, redisKey)
	if err != nil {
		return
	}
	if res.Int() > 0 {
		return
	}
	err = g.Redis().SetEX(c.ctx, redisKey, 1, int64(86400))
	if err != nil {
		return
	}
	var bets []*do.Canada28Bets
	err = dao.Canada28Bets.Ctx(c.ctx).
		Where("period_number=? AND status=?", draw.PeriodNumber, "pending").
		Fields("id", "user_id", "bet_name", "amount", "multiplier").
		Scan(&bets)
	if err != nil {
		g.Log().Error(c.ctx, g.Map{
			"error":         "get pending bets error",
			"msg":           err.Error(),
			"period_number": draw.PeriodNumber,
		})
		return
	}

	// 拼接信息
	message := ""
	if len(bets) >= 0 {
		// 获取所有涉及的用户ID
		userIds := make([]string, 0)
		userIdSet := make(map[string]bool)
		for _, bet := range bets {
			userId := gconv.String(bet.UserId)
			if !userIdSet[userId] {
				userIds = append(userIds, userId)
				userIdSet[userId] = true
			}
		}
		// 一次性查询所有用户信息
		var users []*do.Users
		err = dao.Users.Ctx(c.ctx).Where("uuid IN(?)", userIds).Fields("id", "nickname", "avatar", "balance", "uuid").Scan(&users)
		if err != nil {
			g.Log().Error(c.ctx, g.Map{
				"error": "get users error",
				"msg":   err.Error(),
			})
			return
		}

		// 转换为map方便查询
		userMap := make(map[string]*do.Users)
		for _, user := range users {
			userMap[gconv.String(user.Uuid)] = user
		}
		for _, bet := range bets {
			user := userMap[gconv.String(bet.UserId)]
			if user == nil {
				continue
			}
			message += fmt.Sprintf("%s bet %s, amount %.2f\n", user.Nickname, bet.BetName, gconv.Float64(bet.Amount))
		}

	}
	message += fmt.Sprintf("Period %s is locked! Betting closed 30 seconds before draw.", draw.PeriodNumber)
	err = game_canada28_ws.BroadcastToAllUsers(c.ctx, g.Map{
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
		g.Log().Warning(c.ctx, g.Map{
			"error": "send canada28 alert msg to all users error",
			"msg":   err.Error(),
		})
	}
}
