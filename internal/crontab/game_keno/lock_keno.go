package crontab_game_keno

import (
	"context"
	"demo/internal/dao"
	"demo/internal/model/do"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

type LockKenoCron struct {
	ctx context.Context
}

func NewLockKenoCron(ctx context.Context) *LockKenoCron {
	return &LockKenoCron{
		ctx: ctx,
	}
}

// 提前30秒发送锁定通知
func (c *LockKenoCron) DoLock() error {
	var draw *do.KenoDraws
	// 在开奖时间前30秒发送锁定通知
	lockTime := gtime.Now().Add(30 * time.Second)
	err := dao.KenoDraws.Ctx(c.ctx).Where("status in (0,1) and end_at < ?", lockTime).Order("id desc").Scan(&draw)
	if err != nil {
		return err
	}
	if draw == nil {
		return nil
	}
	return nil
}
