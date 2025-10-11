package cmd

import (
	"context"
	crontab_game_bingo28 "demo/internal/crontab/game_bingo28"
	crontab_game_canada28 "demo/internal/crontab/game_canada28"
	crontab_game_keno "demo/internal/crontab/game_keno"

	"github.com/gogf/gf/v2/os/gcron"
)

type Cron struct{}

func RegisterCron(ctx context.Context) {
	//c := &Cron{}
	// 读取加拿大28开奖结果
	gcron.Add(ctx, "* * * * * *", func(ctx context.Context) {
		crontab_game_canada28.NewSyncCanada28ResCron(ctx).DoSync()
		crontab_game_bingo28.NewSyncBingo28ResCron(ctx).DoSync()
		crontab_game_keno.NewSyncKenoResCron(ctx).DoSync()
	})
	// 锁定加拿大28开奖结果
	gcron.Add(ctx, "* * * * * *", func(ctx context.Context) {
		crontab_game_canada28.NewLockCanada28Cron(ctx).DoLock()
		crontab_game_bingo28.NewLockBingo28Cron(ctx).DoLock()
		crontab_game_keno.NewLockKenoCron(ctx).DoLock()
	})
}
