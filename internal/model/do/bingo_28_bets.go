// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Bingo28Bets is the golang structure of table game_bingo28_bets for DAO operations like Where/Data.
type Bingo28Bets struct {
	g.Meta       `orm:"table:game_bingo28_bets, do:true"`
	Id           interface{} // 主键ID
	MerchantId   interface{} // 商户ID
	UserId       interface{} // 用户ID
	PeriodNumber interface{} // 期号
	BetType      interface{} // 投注类型：high/low/odd/even/num_0等
	BetName      interface{} // 投注名称：High/Low/Number 0等
	Amount       interface{} // 投注金额
	Multiplier   interface{} // 投注时的赔率
	Status       interface{} // 状态：pending-等待开奖，win-已中奖，lose-未中奖，cancel-已取消
	Ip           interface{} // 投注IP地址
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 删除时间
}
