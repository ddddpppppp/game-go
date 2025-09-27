package crontab_game_canada28

import (
	"context"
	consts_sync "demo/internal/consts/sync"
	"demo/internal/dao"
	"demo/internal/model/do"
	"fmt"
	"math"
	"time"

	game_canada28_ws "demo/internal/service/game_canada28_ws"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type SyncCanada28ResCron struct {
	ctx context.Context
}

func NewSyncCanada28ResCron(ctx context.Context) *SyncCanada28ResCron {
	return &SyncCanada28ResCron{
		ctx: ctx,
	}
}

// 1天未完成的直接设为失败
func (c *SyncCanada28ResCron) DoSync() error {
	var draw *do.Canada28Draws
	err := dao.Canada28Draws.Ctx(c.ctx).Where("status in (0,1) and end_at < ?", gtime.Now()).Order("id desc").Scan(&draw)
	if err != nil {
		return err
	}
	if draw == nil {
		return nil
	}
	// 当前时间
	var redisKey = fmt.Sprintf(consts_sync.SyncCanada28ResKey, draw.Id)
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
	dao.Canada28Draws.Ctx(c.ctx).Where("draw_id=?", draw.Id).Update("status=1")

	// 检查开奖时间是否超过1分钟
	timeSinceEnd := gtime.Now().Sub(draw.EndAt)
	timeSinceEndSeconds := int(timeSinceEnd.Seconds())
	useBackupAPI := timeSinceEndSeconds > 60

	var res28 string
	var resMap28 map[string]interface{}
	var list []map[string]interface{}

	// 主API地址
	primaryAPI := "https://apigx.cn/token/63354220977111f08ab63394392e5b2d/code/jnd28/rows/3/type/fastest.json"
	// 备用API地址
	backupAPI := "https://super.pc28998.com/history/JND28"

	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	if useBackupAPI {
		g.Log().Info(c.ctx, g.Map{
			"msg":            "Using backup API due to timeout",
			"period_number":  draw.PeriodNumber,
			"time_since_end": timeSinceEnd.String(),
		})

		// 使用备用API
		res28 = g.Client().SetHeader("User-Agent", userAgent).GetContent(c.ctx, backupAPI)
		resMap28 = gconv.Map(res28)

		// 备用API格式检查：检查是否有code字段且为1
		if gconv.Int(resMap28["code"]) != 1 {
			g.Log().Error(c.ctx, g.Map{
				"error":         "sync canada28 backup API error",
				"msg":           resMap28["msg"],
				"period_number": draw.PeriodNumber,
				"raw":           res28,
			})
			return gerror.New(gconv.String(resMap28["msg"]))
		}

		// 备用API返回的数据格式不同，需要特殊处理
		if resMap28["data"] == nil {
			g.Log().Error(c.ctx, g.Map{
				"error":         "sync canada28 backup API data is null",
				"msg":           "Backup API response data is null",
				"period_number": draw.PeriodNumber,
				"raw":           res28,
			})
			return gerror.New("Backup API response data is null")
		}
		list = gconv.Maps(resMap28["data"])
	} else {
		// 使用主API
		res28 = g.Client().SetHeader("User-Agent", userAgent).GetContent(c.ctx, primaryAPI)
		resMap28 = gconv.Map(res28)

		// 主接口格式检查：检查是否有data字段
		if resMap28["data"] == nil {
			g.Log().Error(c.ctx, g.Map{
				"error":         "sync canada28 primary API error",
				"msg":           "Primary API response data is null",
				"period_number": draw.PeriodNumber,
				"raw":           res28,
			})
			return gerror.New("Primary API response data is null")
		}
		list = gconv.Maps(resMap28["data"])
	}
	for _, item := range list {
		if gconv.String(item["expect"]) == gconv.String(draw.PeriodNumber) {
			var resultNumbers = gstr.Explode(",", gconv.String(item["opencode"]))
			var resultSum = 0
			for _, number := range resultNumbers {
				resultSum += gconv.Int(number)
			}
			draw.ResultNumbers = resultNumbers
			draw.ResultSum = resultSum
			draw.DrawAt = gtime.Now()
			draw.Status = 2
			_, err = dao.Canada28Draws.Ctx(c.ctx).Where("id=?", draw.Id).Update(draw)
			if err != nil {
				g.Log().Error(c.ctx, g.Map{
					"error":         "update canada28 draws error",
					"msg":           err.Error(),
					"period_number": draw.PeriodNumber,
				})
				return err
			}
			// 异步发送开奖通知
			go c.drawResult(draw)
			// 异步发放奖励
			go c.bonus(draw)
			// 插入下一期期数
			currentTime := draw.EndAt.Time
			// 如果当前时间是10:56:30，下一期结束时间是11:33:00
			var nextEndAt time.Time
			if currentTime.Hour() == 10 && currentTime.Minute() == 56 && currentTime.Second() == 30 {
				nextEndAt = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 11, 33, 0, 0, currentTime.Location())
			} else {
				// 其他情况下+210秒
				nextEndAt = currentTime.Add(time.Second * 210)
			}
			_, err = dao.Canada28Draws.Ctx(c.ctx).Insert(g.Map{
				"period_number": gconv.Int(draw.PeriodNumber) + 1,
				"status":        0,
				"start_at":      draw.EndAt,
				"end_at":        nextEndAt,
				"draw_at":       nextEndAt,
			})
			if err != nil {
				g.Log().Error(c.ctx, g.Map{
					"error":         "insert canada28 draws error",
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

// 计算动态赔率
func (c *SyncCanada28ResCron) calculateDynamicOdds(betType string, originalOdds float64, resultSum int, merchantId string) float64 {
	// 查询适用的动态赔率规则
	var rules []*do.Canada28DynamicOdds
	err := dao.Canada28DynamicOdds.Ctx(c.ctx).
		Where("merchant_id = ? AND status = 1", merchantId).
		Order("priority DESC").
		Scan(&rules)

	if err != nil {
		g.Log().Error(c.ctx, g.Map{
			"error": "query dynamic odds rules failed",
			"msg":   err.Error(),
		})
		return originalOdds
	}

	// 检查每个规则是否适用
	for _, rule := range rules {
		if c.checkRuleCondition(rule, resultSum) {
			// 解析赔率调整配置
			adjustmentsJson := gconv.String(rule.BetTypeAdjustments)
			if adjustmentsJson != "" {
				json := gjson.New(adjustmentsJson)
				if newOdds := json.Get(betType); newOdds != nil {
					return newOdds.Float64()
				}
			}
		}
	}

	return originalOdds
}

// 检查规则条件是否满足
func (c *SyncCanada28ResCron) checkRuleCondition(rule *do.Canada28DynamicOdds, resultSum int) bool {
	triggerCondition := gconv.String(rule.TriggerCondition)
	triggerValues := gconv.String(rule.TriggerValues)

	switch triggerCondition {
	case "sum_in":
		json := gjson.New(triggerValues)
		values := json.Array()
		for _, v := range values {
			if gconv.Int(v) == resultSum {
				return true
			}
		}
		return false
	case "sum_range":
		json := gjson.New(triggerValues)
		min := gconv.Int(json.Get("min"))
		max := gconv.Int(json.Get("max"))
		return resultSum >= min && resultSum <= max
	case "sum_exact":
		exact := gconv.Int(gjson.New(triggerValues))
		return resultSum == exact
	default:
		return false
	}
}

func (c *SyncCanada28ResCron) bonus(draw *do.Canada28Draws) {
	// 获取该期所有等待开奖的投注记录
	var redisKey = fmt.Sprintf(consts_sync.SyncCanada28BonusKey, draw.Id)
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
		Fields("id", "user_id", "bet_name", "bet_type", "amount", "multiplier", "merchant_id").
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
		return
	}

	// 计算开奖结果相关数据
	resultSum := gconv.Int(draw.ResultSum)
	isOdd := resultSum%2 == 1
	isEven := !isOdd
	isHigh := resultSum >= 14 && resultSum <= 27
	isLow := resultSum >= 0 && resultSum <= 13
	isExtremeHigh := resultSum >= 22 && resultSum <= 27
	isExtremeLow := resultSum >= 0 && resultSum <= 5

	// 检查特殊投注类型
	numbers := gconv.Strings(draw.ResultNumbers)
	isTriple := len(numbers) == 3 && numbers[0] == numbers[1] && numbers[1] == numbers[2]
	isPair := len(numbers) == 3 && (numbers[0] == numbers[1] || numbers[1] == numbers[2] || numbers[0] == numbers[2])
	isStraight := false
	if len(numbers) == 3 {
		nums := []int{gconv.Int(numbers[0]), gconv.Int(numbers[1]), gconv.Int(numbers[2])}
		// 检查是否为连续数字（任意顺序）
		for i := 0; i < 3; i++ {
			for j := i + 1; j < 3; j++ {
				if nums[i] > nums[j] {
					nums[i], nums[j] = nums[j], nums[i]
				}
			}
		}
		isStraight = nums[1] == nums[0]+1 && nums[2] == nums[1]+1
	}

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
	err = dao.Users.Ctx(c.ctx).Where("uuid IN(?)", userIds).Fields("id", "nickname", "avatar", "balance", "uuid", "balance_frozen").Scan(&users)
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

	// 批量收集需要更新的数据
	type betUpdate struct {
		betId      int64
		status     string
		multiplier *float64 // 可选字段，仅在需要更新时使用
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
	userBalanceMap := make(map[string]float64)       // 用于累积每个用户的余额变化
	userFrozenBalanceMap := make(map[string]float64) // 用于累积每个用户的冻结余额变化

	message := ""

	// 批量查询投注相关的余额记录
	betIds := make([]string, 0, len(bets))
	for _, bet := range bets {
		betIds = append(betIds, gconv.String(bet.Id))
	}

	// 查询所有投注对应的balance记录（投注时的扣款记录）
	var balanceRecords []*do.UserBalances
	var frozenBalanceRecords []*do.UserFrozenBalances

	if len(betIds) > 0 {
		// 查询普通余额记录 (type=game_bet)
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

		// 查询冻结余额记录 (type=game_bet)
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

	// 构建投注ID到余额记录的映射
	betToBalanceMap := make(map[string]*do.UserBalances)
	betToFrozenBalanceMap := make(map[string]*do.UserFrozenBalances)

	for _, record := range balanceRecords {
		betToBalanceMap[gconv.String(record.RelatedId)] = record
	}

	for _, record := range frozenBalanceRecords {
		betToFrozenBalanceMap[gconv.String(record.RelatedId)] = record
	}

	// 遍历所有投注记录计算中奖情况
	for _, bet := range bets {
		betType := gconv.String(bet.BetType)
		betAmount := gconv.Float64(bet.Amount)
		multiplier := gconv.Float64(bet.Multiplier)
		userId := gconv.String(bet.UserId)
		betId := gconv.Int64(bet.Id)
		betIdStr := gconv.String(bet.Id)
		merchantId := gconv.String(bet.MerchantId)

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
			// 如果只有冻结余额记录，说明全部来自冻结余额
			frozenRatio = 1.0
		}
		// 如果只有普通余额记录或都没有，frozenRatio保持为0.0

		isWin := false

		// 根据投注类型判断是否中奖
		switch betType {
		case "high":
			isWin = isHigh
		case "low":
			isWin = isLow
		case "odd":
			isWin = isOdd
		case "even":
			isWin = isEven
		case "extreme_high":
			isWin = isExtremeHigh
		case "extreme_low":
			isWin = isExtremeLow
		case "high_odd":
			isWin = isHigh && isOdd
		case "high_even":
			isWin = isHigh && isEven
		case "low_odd":
			isWin = isLow && isOdd
		case "low_even":
			isWin = isLow && isEven
		case "triple":
			isWin = isTriple
		case "pair":
			isWin = isPair && !isTriple // 对子不包括三条
		case "straight":
			isWin = isStraight
		default:
			// 特码投注 (sum_0 到 sum_27)
			if gstr.HasPrefix(betType, "sum_") {
				targetNum := gconv.Int(gstr.SubStr(betType, 4))
				isWin = resultSum == targetNum
			}
		}

		// 更新投注状态
		if isWin {
			// 计算动态赔率
			finalMultiplier := c.calculateDynamicOdds(betType, multiplier, resultSum, merchantId)

			// 创建投注更新记录
			betUpdateRecord := betUpdate{
				betId:  betId,
				status: "win",
			}

			// 如果动态赔率与原始赔率不同，也要更新multiplier
			if finalMultiplier != multiplier {
				betUpdateRecord.multiplier = &finalMultiplier
			}

			betUpdates = append(betUpdates, betUpdateRecord)

			// 计算奖金（投注金额 * 动态赔率）
			winAmount := betAmount * finalMultiplier

			// 获取用户当前余额（包含之前的累积）
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
			// 冻结余额根据比例获得额外奖金
			frozenWinAmount := winAmount * frozenRatio

			// 处理普通余额部分 - 总是获得全部奖金
			oldBalance := gconv.Float64(user.Balance) + userBalanceMap[userId]
			newBalance := oldBalance + normalWinAmount

			// 累积用户余额变化
			userBalanceMap[userId] += normalWinAmount

			// 每个中奖投注对应一条余额记录
			balanceUpdates = append(balanceUpdates, userBalanceUpdate{
				userUuid:    gconv.String(user.Uuid),
				userId:      gconv.Int(user.Id),
				winAmount:   normalWinAmount,
				description: fmt.Sprintf("Keno Win - Period:%s, Bet Type:%s, Multiplier:%.2f (Original:%.2f)", draw.PeriodNumber, betType, finalMultiplier, multiplier),
				relatedId:   betId,
				oldBalance:  oldBalance,
				newBalance:  newBalance,
			})

			// 处理冻结余额部分 - 根据比例获得额外奖金
			if frozenWinAmount > 0 {
				oldFrozenBalance := gconv.Float64(user.BalanceFrozen) + userFrozenBalanceMap[userId]
				newFrozenBalance := oldFrozenBalance + frozenWinAmount

				// 累积用户冻结余额变化
				userFrozenBalanceMap[userId] += frozenWinAmount

				// 每个中奖投注对应一条冻结余额记录
				frozenBalanceUpdates = append(frozenBalanceUpdates, userFrozenBalanceUpdate{
					userUuid:    gconv.String(user.Uuid),
					userId:      gconv.Int(user.Id),
					winAmount:   frozenWinAmount,
					description: fmt.Sprintf("Keno Win - Period:%s, Bet Type:%s, Multiplier:%.2f (Original:%.2f) [Frozen Balance]", draw.PeriodNumber, betType, finalMultiplier, multiplier),
					relatedId:   betId,
					oldBalance:  oldFrozenBalance,
					newBalance:  newFrozenBalance,
				})
			}

		} else {
			betUpdates = append(betUpdates, betUpdate{
				betId:  betId,
				status: "lose",
			})
		}
	}

	// 开始事务批量更新
	err = g.DB().Transaction(c.ctx, func(ctx context.Context, tx gdb.TX) error {
		// 批量更新投注状态
		for _, update := range betUpdates {
			updateData := g.Map{
				"status":     update.status,
				"updated_at": gtime.Now(),
			}

			// 如果需要更新multiplier
			if update.multiplier != nil {
				updateData["multiplier"] = *update.multiplier
			}

			_, err := tx.Model("game_canada28_bets").
				Where("id=?", update.betId).
				Update(updateData)
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

		dao.Canada28Draws.Ctx(ctx).Where("id=?", draw.Id).Update("status=3")

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
	// 根据userBalanceMap按总赢钱金额排序，取前20名
	type userWin struct {
		userId    string
		nickname  string
		winAmount float64
	}

	var winners []userWin
	for userId, totalWin := range userBalanceMap {
		if totalWin > 0 {
			user := userMap[userId]
			if user != nil {
				winners = append(winners, userWin{
					userId:    userId,
					nickname:  gconv.String(user.Nickname),
					winAmount: totalWin,
				})
			}
		}
	}

	// 按赢钱金额降序排序
	for i := 0; i < len(winners); i++ {
		for j := i + 1; j < len(winners); j++ {
			if winners[j].winAmount > winners[i].winAmount {
				winners[i], winners[j] = winners[j], winners[i]
			}
		}
	}

	// 取前20名
	maxWinners := 20
	if len(winners) > maxWinners {
		winners = winners[:maxWinners]
	}

	// 生成消息
	if len(winners) == 0 {
		message = "What a pity, this period has no winners!"
	} else {
		message = "🎉 Period " + gconv.String(draw.PeriodNumber) + " winner List:\n"
		for i, winner := range winners {
			if winner.winAmount <= 0 {
				continue
			}
			message += fmt.Sprintf("%d. %s won %.2f\n", i+1, winner.nickname, winner.winAmount)
		}
		message += "Congratulations to the top 20 winners!"
	}
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
			"error": "send canada28 res msg to all users error",
			"msg":   err.Error(),
		})
	}

	g.Log().Info(c.ctx, g.Map{
		"msg":                 "bonus distribution completed",
		"period_number":       draw.PeriodNumber,
		"total_bets":          len(bets),
		"balance_wins":        len(userBalanceMap),
		"frozen_balance_wins": len(userFrozenBalanceMap),
	})
}

func (c *SyncCanada28ResCron) drawResult(draw *do.Canada28Draws) {
	err := game_canada28_ws.BroadcastToAllUsers(c.ctx, g.Map{
		"action": "draw_result",
		"data": g.Map{
			"period_number":  draw.PeriodNumber,
			"result_numbers": draw.ResultNumbers,
			"result_sum":     draw.ResultSum,
		},
	})
	if err != nil {
		g.Log().Warning(c.ctx, g.Map{
			"error": "send canada28 res msg to all users error",
			"msg":   err.Error(),
		})
	}
}
