// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Bingo28Draws is the golang structure of table game_bingo28_draws for DAO operations like Where/Data.
type Bingo28Draws struct {
	g.Meta        `orm:"table:game_bingo28_draws, do:true"`
	Id            interface{} // 主键ID
	PeriodNumber  interface{} // 期号，如：3333197
	Status        interface{} // 状态：0-等待开奖，1-开奖中，2-已开奖，3-已结算
	StartAt       *gtime.Time // 开始投注时间
	EndAt         *gtime.Time // 停止投注时间
	DrawAt        *gtime.Time // 开奖时间
	ResultNumbers interface{} // 开奖号码，JSON格式存储三个数字
	ResultSum     interface{} // 开奖结果总和(0-27)
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	DeletedAt     *gtime.Time // 删除时间
}
