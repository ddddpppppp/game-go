package crontab_game_bingo

import (
	"context"
	"demo/internal/dao"
	"demo/internal/model/do"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

type LockBingoCron struct {
	ctx context.Context
}

func NewLockBingoCron(ctx context.Context) *LockBingoCron {
	return &LockBingoCron{
		ctx: ctx,
	}
}

// 提前30秒发送锁定通知
func (c *LockBingoCron) DoLock() error {
	var draw *do.BingoDraws
	// 在开奖时间前30秒发送锁定通知
	lockTime := gtime.Now().Add(30 * time.Second)
	err := dao.BingoDraws.Ctx(c.ctx).Where("status in (0,1) and end_at < ?", lockTime).Order("id desc").Scan(&draw)
	if err != nil {
		return err
	}
	if draw == nil {
		return nil
	}
	return nil
}
