// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BingoDraws is the golang structure of table game_bingo_draws for DAO operations like Where/Data.
type BingoDraws struct {
	g.Meta        `orm:"table:game_bingo_draws, do:true"`
	Id            interface{} // 主键ID
	PeriodNumber  interface{} // 期号，如：3333197
	Status        interface{} // 状态：0-等待开奖，1-开奖中，2-已开奖，3-已结算
	StartAt       *gtime.Time // 开始投注时间
	EndAt         *gtime.Time // 停止投注时间
	DrawAt        *gtime.Time // 开奖时间
	ResultNumbers interface{} // 开奖号码，JSON格式存储三个数字
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	DeletedAt     *gtime.Time // 删除时间
}
