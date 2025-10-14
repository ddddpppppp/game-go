// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BingoBets is the golang structure of table game_bingo_bets for DAO operations like Where/Data.
type BingoBets struct {
	g.Meta          `orm:"table:game_bingo_bets, do:true"`
	Id              interface{} // 主键ID
	MerchantId      interface{} // 商户ID
	UserId          interface{} // 用户ID
	PeriodNumber    interface{} // 期号
	SelectedNumbers interface{} // 玩家选择的10个号码(1-80)，JSON格式: [1,5,12,...]
	DrawnNumbers    interface{} // 开出的20个号码(1-80)，JSON格式: [3,7,12,...]
	MatchedNumbers  interface{} // 匹配的号码，JSON格式: [12,...]
	MatchCount      interface{} // 匹配数量：0-10
	Amount          interface{} // 投注金额
	Multiplier      interface{} // 投注时的赔率（基于匹配数量）
	WinAmount       interface{} // 中奖金额
	Status          interface{} // 状态：pending-等待开奖，win-已中奖，lose-未中奖，cancel-已取消
	Ip              interface{} // 投注IP地址
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
	SettledAt       *gtime.Time // 结算时间
	DeletedAt       *gtime.Time // 删除时间
}
