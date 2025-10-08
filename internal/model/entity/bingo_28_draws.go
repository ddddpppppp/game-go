// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Bingo28Draws is the golang structure for table bingo28_draws.
type Bingo28Draws struct {
	Id            uint64      `json:"id"            orm:"id"             description:"主键ID"`                        // 主键ID
	PeriodNumber  string      `json:"periodNumber"  orm:"period_number"  description:"期号，如：3333197"`                // 期号，如：3333197
	Status        int         `json:"status"        orm:"status"         description:"状态：0-等待开奖，1-开奖中，2-已开奖，3-已结算"` // 状态：0-等待开奖，1-开奖中，2-已开奖，3-已结算
	StartAt       *gtime.Time `json:"startAt"       orm:"start_at"       description:"开始投注时间"`                      // 开始投注时间
	EndAt         *gtime.Time `json:"endAt"         orm:"end_at"         description:"停止投注时间"`                      // 停止投注时间
	DrawAt        *gtime.Time `json:"drawAt"        orm:"draw_at"        description:"开奖时间"`                        // 开奖时间
	ResultNumbers string      `json:"resultNumbers" orm:"result_numbers" description:"开奖号码，JSON格式存储三个数字"`           // 开奖号码，JSON格式存储三个数字
	ResultSum     int         `json:"resultSum"     orm:"result_sum"     description:"开奖结果总和(0-27)"`                // 开奖结果总和(0-27)
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:"创建时间"`                        // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:"更新时间"`                        // 更新时间
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"     description:"删除时间"`                        // 删除时间
}
