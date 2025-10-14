package crontab_game_bingo

import (
	"context"
	consts_sync "demo/internal/consts/sync"
	"demo/internal/dao"
	"demo/internal/model/do"
	"fmt"
	"math"
	"time"

	game_bingo_ws "demo/internal/service/game_bingo_ws"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type SyncBingoResCron struct {
	ctx context.Context
}

func NewSyncBingoResCron(ctx context.Context) *SyncBingoResCron {
	return &SyncBingoResCron{
		ctx: ctx,
	}
}

// Bingo开奖同步 (从API获取20个号码 1-80)
func (c *SyncBingoResCron) DoSync() error {
	var draw *do.BingoDraws
	err := dao.BingoDraws.Ctx(c.ctx).Where("status in (0,1) and end_at < ?", gtime.Now()).Order("id desc").Scan(&draw)
	if err != nil {
		return err
	}
	if draw == nil {
		return nil
	}

	// Redis锁防止并发开奖
	var redisKey = fmt.Sprintf(consts_sync.SyncBingoResKey, draw.Id)
	res, err := g.Redis().Get(c.ctx, redisKey)
	if err != nil {
		return err
	}
	if res.Int() > 0 {
		return nil
	}
	err = g.Redis().SetEX(c.ctx, redisKey, 1, int64(5))
	if err != nil {
		return err
	}

	// 标记为开奖中
	dao.BingoDraws.Ctx(c.ctx).Where("id=?", draw.Id).Update("status=1")

	// 检查开奖时间是否超过1分钟，决定使用主API还是备用API
	timeSinceEnd := gtime.Now().Sub(draw.EndAt)
	timeSinceEndSeconds := int(timeSinceEnd.Seconds())
	useBackupAPI := timeSinceEndSeconds > 60

	var resBingo string
	var resMapBingo map[string]interface{}
	var list []map[string]interface{}

	// 主API地址 - Bingo接口 (返回20个号码 1-80)
	primaryAPI := "https://apigx.cn/token/92995bbea69911f0a45553dcc3f98e8b/code/twbg/rows/3.json"
	// 备用API地址 - Bingo接口
	backupAPI := "https://vip.apigx.cn:2096/token/92995bbea69911f0a45553dcc3f98e8b/code/twbg/rows/3.json"

	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	if useBackupAPI {
		g.Log().Info(c.ctx, g.Map{
			"msg":            "Using backup Bingo API due to timeout",
			"period_number":  draw.PeriodNumber,
			"time_since_end": timeSinceEnd.String(),
		})
		resBingo = g.Client().SetHeader("User-Agent", userAgent).GetContent(c.ctx, backupAPI)
	} else {
		resBingo = g.Client().SetHeader("User-Agent", userAgent).GetContent(c.ctx, primaryAPI)
	}

	resMapBingo = gconv.Map(resBingo)

	// API格式检查
	if resMapBingo["data"] == nil {
		g.Log().Error(c.ctx, g.Map{
			"error":         "sync bingo API error",
			"msg":           "API response data is null",
			"period_number": draw.PeriodNumber,
			"raw":           resBingo,
		})
		return fmt.Errorf("API response data is null")
	}
	list = gconv.Maps(resMapBingo["data"])

	// 查找匹配的期号
	for _, item := range list {
		if gconv.String(item["expect"]) == gconv.String(draw.PeriodNumber) {
			// 获取20个开奖号码 (opencode格式: "2,66,49,6,58,26,71,44,30,56,33,23,31,25,8,79,74,12,63,57")
			var resultNumbersStr = gstr.Explode(",", gconv.String(item["opencode"]))
			resultNumbers := make([]int, 0, 20)

			// 转换为整数并排序
			for _, numStr := range resultNumbersStr {
				num := gconv.Int(numStr)
				if num >= 1 && num <= 80 {
					resultNumbers = append(resultNumbers, num)
				}
			}

			// 排序
			for i := 0; i < len(resultNumbers); i++ {
				for j := i + 1; j < len(resultNumbers); j++ {
					if resultNumbers[i] > resultNumbers[j] {
						resultNumbers[i], resultNumbers[j] = resultNumbers[j], resultNumbers[i]
					}
				}
			}

			// 更新开奖结果
			draw.ResultNumbers = resultNumbers
			draw.DrawAt = gtime.Now()
			draw.Status = 2

			_, err = dao.BingoDraws.Ctx(c.ctx).Where("id=?", draw.Id).Update(g.Map{
				"result_numbers": resultNumbers,
				"draw_at":        gtime.Now(),
				"status":         2,
			})
			if err != nil {
				g.Log().Error(c.ctx, g.Map{
					"error":         "update bingo draws error",
					"msg":           err.Error(),
					"period_number": draw.PeriodNumber,
				})
				return err
			}

			g.Log().Info(c.ctx, g.Map{
				"msg":            "Bingo draw completed",
				"period_number":  draw.PeriodNumber,
				"result_numbers": resultNumbers,
			})

			// 异步发送开奖通知 (WebSocket)
			go c.drawResult(draw)

			// 异步结算奖励
			go c.bonus(draw)

			// 创建下一期
			currentTime := draw.EndAt.Time
			// 如果当前时间是15:55:00，下一期结束时间是23:05:00
			var nextEndAt time.Time
			if currentTime.Hour() == 15 && currentTime.Minute() == 55 && currentTime.Second() == 0 {
				nextEndAt = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 05, 0, 0, currentTime.Location())
			} else {
				// 其他情况下+300秒
				nextEndAt = currentTime.Add(time.Second * 300)
			}

			_, err = dao.BingoDraws.Ctx(c.ctx).Insert(g.Map{
				"period_number": gconv.Int(draw.PeriodNumber) + 1,
				"status":        0,
				"start_at":      draw.EndAt,
				"end_at":        nextEndAt,
				"draw_at":       nextEndAt,
			})
			if err != nil {
				g.Log().Error(c.ctx, g.Map{
					"error":         "insert next bingo period error",
					"msg":           err.Error(),
					"period_number": draw.PeriodNumber,
				})
				return err
			}

			break
		}
	}

	return nil
}

// 结算奖金 (Bingo BCLC规则)
func (c *SyncBingoResCron) bonus(draw *do.BingoDraws) {
	// Redis锁防止重复结算
	var redisKey = fmt.Sprintf(consts_sync.SyncBingoBonusKey, draw.Id)
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

	// 获取该期所有待结算的投注
	var bets []*do.BingoBets
	err = dao.BingoBets.Ctx(c.ctx).
		Where("period_number=? AND status=?", draw.PeriodNumber, "pending").
		Fields("id", "user_id", "selected_numbers", "amount", "merchant_id").
		Scan(&bets)
	if err != nil {
		g.Log().Error(c.ctx, g.Map{
			"error":         "get pending bets error",
			"msg":           err.Error(),
			"period_number": draw.PeriodNumber,
		})
		return
	}

	if len(bets) == 0 {
		g.Log().Info(c.ctx, g.Map{
			"msg":           "no bets to settle",
			"period_number": draw.PeriodNumber,
		})
		return
	}

	// 获取开出的号码
	drawnNumbers := gconv.Ints(draw.ResultNumbers)

	// 批量查询投注相关的余额记录
	betIds := make([]string, 0, len(bets))
	for _, bet := range bets {
		betIds = append(betIds, gconv.String(bet.Id))
	}

	var balanceRecords []*do.UserBalances
	var frozenBalanceRecords []*do.UserFrozenBalances

	if len(betIds) > 0 {
		err = dao.UserBalances.Ctx(c.ctx).
			Where("related_id IN(?) AND type=?", betIds, "game_bet").
			Fields("id", "user_id", "amount", "related_id").
			Scan(&balanceRecords)
		if err != nil {
			g.Log().Error(c.ctx, g.Map{
				"error": "get balance records error",
				"msg":   err.Error(),
			})
			return
		}

		err = dao.UserFrozenBalances.Ctx(c.ctx).
			Where("related_id IN(?) AND type=?", betIds, "game_bet").
			Fields("id", "user_id", "amount", "related_id").
			Scan(&frozenBalanceRecords)
		if err != nil {
			g.Log().Error(c.ctx, g.Map{
				"error": "get frozen balance records error",
				"msg":   err.Error(),
			})
			return
		}
	}

	betToBalanceMap := make(map[string]*do.UserBalances)
	betToFrozenBalanceMap := make(map[string]*do.UserFrozenBalances)

	for _, record := range balanceRecords {
		betToBalanceMap[gconv.String(record.RelatedId)] = record
	}

	for _, record := range frozenBalanceRecords {
		betToFrozenBalanceMap[gconv.String(record.RelatedId)] = record
	}

	// 批量收集需要更新的数据
	type betUpdate struct {
		betId          int64
		status         string
		drawnNumbers   []int
		matchedNumbers []int
		matchCount     int
		multiplier     float64
		winAmount      float64
	}

	type userBalanceUpdate struct {
		userId      int
		userUuid    string
		winAmount   float64
		description string
		relatedId   int64
		oldBalance  float64
		newBalance  float64
	}

	type userFrozenBalanceUpdate struct {
		userId      int
		userUuid    string
		winAmount   float64
		description string
		relatedId   int64
		oldBalance  float64
		newBalance  float64
	}

	var betUpdates []betUpdate
	var balanceUpdates []userBalanceUpdate
	var frozenBalanceUpdates []userFrozenBalanceUpdate
	userBalanceMap := make(map[string]float64)
	userFrozenBalanceMap := make(map[string]float64)

	// 查询赔率表
	var betTypes []*do.BingoBetTypes
	err = dao.BingoBetTypes.Ctx(c.ctx).
		Where("merchant_id=? AND status=1", bets[0].MerchantId).
		Fields("id", "type_key", "type_name", "odds").
		Scan(&betTypes)
	if err != nil {
		g.Log().Error(c.ctx, g.Map{
			"error": "get bet types error",
			"msg":   err.Error(),
		})
		return
	}

	// 构建赔率映射 (match_0 -> odds)
	oddsMap := make(map[int]float64)
	for _, betType := range betTypes {
		typeKey := gconv.String(betType.TypeKey)
		if gstr.HasPrefix(typeKey, "match_") {
			matchCount := gconv.Int(gstr.SubStr(typeKey, 6))
			oddsMap[matchCount] = gconv.Float64(betType.Odds)
		}
	}

	// 查询所有涉及的用户
	userIds := make([]string, 0)
	userIdSet := make(map[string]bool)
	for _, bet := range bets {
		userId := gconv.String(bet.UserId)
		if !userIdSet[userId] {
			userIds = append(userIds, userId)
			userIdSet[userId] = true
		}
	}

	var users []*do.Users
	err = dao.Users.Ctx(c.ctx).Where("uuid IN(?)", userIds).
		Fields("id", "nickname", "avatar", "balance", "uuid", "balance_frozen").
		Scan(&users)
	if err != nil {
		g.Log().Error(c.ctx, g.Map{
			"error": "get users error",
			"msg":   err.Error(),
		})
		return
	}

	userMap := make(map[string]*do.Users)
	for _, user := range users {
		userMap[gconv.String(user.Uuid)] = user
	}

	// 遍历所有投注记录并计算匹配
	for _, bet := range bets {
		betAmount := gconv.Float64(bet.Amount)
		userId := gconv.String(bet.UserId)
		betId := gconv.Int64(bet.Id)
		betIdStr := gconv.String(bet.Id)

		// 解析玩家选择的号码
		selectedNumbersJson := gconv.String(bet.SelectedNumbers)
		selectedNumbers := gconv.Ints(selectedNumbersJson)

		// 计算匹配的号码
		matchedNumbers := make([]int, 0)
		for _, num := range selectedNumbers {
			for _, drawn := range drawnNumbers {
				if num == drawn {
					matchedNumbers = append(matchedNumbers, num)
					break
				}
			}
		}
		matchCount := len(matchedNumbers)

		// 获取对应的赔率
		multiplier := oddsMap[matchCount]
		winAmount := betAmount * multiplier

		// 计算该投注的冻结余额占比
		var frozenRatio float64 = 0.0
		balanceRecord := betToBalanceMap[betIdStr]
		frozenBalanceRecord := betToFrozenBalanceMap[betIdStr]
		if balanceRecord != nil && frozenBalanceRecord != nil {
			totalAmount := math.Abs(gconv.Float64(balanceRecord.Amount))
			frozenAmount := math.Abs(gconv.Float64(frozenBalanceRecord.Amount))
			if totalAmount > 0 {
				frozenRatio = frozenAmount / totalAmount
			}
		} else if frozenBalanceRecord != nil {
			frozenRatio = 1.0
		}

		// 创建投注更新记录
		betUpdateRecord := betUpdate{
			betId: betId,
			status: func() string {
				if winAmount > 0 {
					return "win"
				}
				return "lose"
			}(),
			drawnNumbers:   drawnNumbers,
			matchedNumbers: matchedNumbers,
			matchCount:     matchCount,
			multiplier:     multiplier,
			winAmount:      winAmount,
		}
		betUpdates = append(betUpdates, betUpdateRecord)

		// 如果中奖，处理余额
		if winAmount > 0 {
			user := userMap[userId]
			if user == nil {
				g.Log().Error(c.ctx, g.Map{
					"error":   "user not found in map",
					"user_id": userId,
				})
				continue
			}

			// 普通余额获得全部奖金
			normalWinAmount := winAmount
			frozenWinAmount := winAmount * frozenRatio

			// 处理普通余额
			oldBalance := gconv.Float64(user.Balance) + userBalanceMap[userId]
			newBalance := oldBalance + normalWinAmount
			userBalanceMap[userId] += normalWinAmount

			balanceUpdates = append(balanceUpdates, userBalanceUpdate{
				userUuid:    gconv.String(user.Uuid),
				userId:      gconv.Int(user.Id),
				winAmount:   normalWinAmount,
				description: fmt.Sprintf("Bingo Win - Period:%s, Matches:%d, Odds:%.2f", draw.PeriodNumber, matchCount, multiplier),
				relatedId:   betId,
				oldBalance:  oldBalance,
				newBalance:  newBalance,
			})

			// 处理冻结余额
			if frozenWinAmount > 0 {
				oldFrozenBalance := gconv.Float64(user.BalanceFrozen) + userFrozenBalanceMap[userId]
				newFrozenBalance := oldFrozenBalance + frozenWinAmount
				userFrozenBalanceMap[userId] += frozenWinAmount

				frozenBalanceUpdates = append(frozenBalanceUpdates, userFrozenBalanceUpdate{
					userUuid:    gconv.String(user.Uuid),
					userId:      gconv.Int(user.Id),
					winAmount:   frozenWinAmount,
					description: fmt.Sprintf("Bingo Win - Period:%s, Matches:%d, Odds:%.2f [Frozen]", draw.PeriodNumber, matchCount, multiplier),
					relatedId:   betId,
					oldBalance:  oldFrozenBalance,
					newBalance:  newFrozenBalance,
				})
			}
		}
	}

	// 开始事务批量更新
	err = g.DB().Transaction(c.ctx, func(ctx context.Context, tx gdb.TX) error {
		// 批量更新投注状态和结果
		for _, update := range betUpdates {
			_, err := tx.Model("game_bingo_bets").
				Where("id=?", update.betId).
				Update(g.Map{
					"drawn_numbers":   update.drawnNumbers,
					"matched_numbers": update.matchedNumbers,
					"match_count":     update.matchCount,
					"multiplier":      update.multiplier,
					"win_amount":      update.winAmount,
					"status":          update.status,
					"settled_at":      gtime.Now(),
					"updated_at":      gtime.Now(),
				})
			if err != nil {
				return err
			}
		}

		// 批量更新用户余额
		for userId, totalChange := range userBalanceMap {
			if totalChange <= 0 {
				continue
			}

			// 获取用户当前余额
			user := userMap[userId]
			if user == nil {
				g.Log().Error(ctx, g.Map{
					"error":   "user not found in map",
					"user_id": userId,
				})
				continue
			}

			oldBalance := gconv.Float64(user.Balance)
			newBalance := oldBalance + totalChange

			// 更新用户余额
			_, err = tx.Model("game_users").
				Where("uuid=?", userId).
				Update(g.Map{"balance": newBalance, "updated_at": gtime.Now()})
			if err != nil {
				return err
			}
		}

		// 批量更新用户冻结余额
		for userId, totalChange := range userFrozenBalanceMap {
			if totalChange <= 0 {
				continue
			}

			// 获取用户当前冻结余额
			user := userMap[userId]
			if user == nil {
				g.Log().Error(ctx, g.Map{
					"error":   "user not found in map",
					"user_id": userId,
				})
				continue
			}

			oldFrozenBalance := gconv.Float64(user.BalanceFrozen)
			newFrozenBalance := oldFrozenBalance + totalChange

			// 更新用户冻结余额
			_, err = tx.Model("game_users").
				Where("uuid=?", userId).
				Update(g.Map{"balance_frozen": newFrozenBalance, "updated_at": gtime.Now()})
			if err != nil {
				return err
			}
		}

		// 批量插入余额变动记录 - 每个中奖投注对应一条记录
		for _, balanceUpdate := range balanceUpdates {
			_, err = tx.Model("game_user_balances").Insert(g.Map{
				"user_id":        balanceUpdate.userId,
				"type":           "game_win",
				"amount":         balanceUpdate.winAmount,
				"balance_before": balanceUpdate.oldBalance,
				"balance_after":  balanceUpdate.newBalance,
				"description":    balanceUpdate.description,
				"related_id":     balanceUpdate.relatedId,
				"created_at":     gtime.Now(),
			})
			if err != nil {
				return err
			}
		}

		// 批量插入冻结余额变动记录 - 每个中奖投注对应一条记录
		for _, frozenBalanceUpdate := range frozenBalanceUpdates {
			_, err = tx.Model("game_user_frozen_balances").Insert(g.Map{
				"user_id":        frozenBalanceUpdate.userId,
				"type":           "game_win",
				"amount":         frozenBalanceUpdate.winAmount,
				"balance_before": frozenBalanceUpdate.oldBalance,
				"balance_after":  frozenBalanceUpdate.newBalance,
				"description":    frozenBalanceUpdate.description,
				"related_id":     frozenBalanceUpdate.relatedId,
				"created_at":     gtime.Now(),
			})
			if err != nil {
				return err
			}
		}

		dao.BingoDraws.Ctx(ctx).Where("id=?", draw.Id).Update("status=3")

		return err
	})

	if err != nil {
		g.Log().Error(c.ctx, g.Map{
			"error":         "bonus distribution error",
			"msg":           err.Error(),
			"period_number": draw.PeriodNumber,
		})
		return
	}

	// 统计结算信息
	totalWinAmount := 0.0
	winCount := 0
	for _, update := range betUpdates {
		if update.winAmount > 0 {
			winCount++
			totalWinAmount += update.winAmount
		}
	}

	g.Log().Info(c.ctx, g.Map{
		"msg":           "Bingo settlement completed",
		"period_number": draw.PeriodNumber,
		"total_bets":    len(bets),
		"win_bets":      winCount,
		"total_payout":  totalWinAmount,
		"drawn_numbers": drawnNumbers,
	})
}

// 推送开奖结果到 WebSocket (只发开奖结果，不发群聊消息)
func (c *SyncBingoResCron) drawResult(draw *do.BingoDraws) {
	err := game_bingo_ws.BroadcastToAllUsers(c.ctx, g.Map{
		"action": "draw_result",
		"data": g.Map{
			"period_number":  draw.PeriodNumber,
			"result_numbers": draw.ResultNumbers,
			"draw_at":        draw.DrawAt,
		},
	})
	if err != nil {
		g.Log().Warning(c.ctx, g.Map{
			"error":         "send bingo draw result error",
			"msg":           err.Error(),
			"period_number": draw.PeriodNumber,
		})
	} else {
		g.Log().Info(c.ctx, g.Map{
			"msg":            "Bingo draw result broadcasted",
			"period_number":  draw.PeriodNumber,
			"result_numbers": draw.ResultNumbers,
		})
	}
}
