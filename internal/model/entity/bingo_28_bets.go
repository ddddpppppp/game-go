// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Bingo28Bets is the golang structure for table bingo28_bets.
type Bingo28Bets struct {
	Id           uint64      `json:"id"           orm:"id"            description:"主键ID"`                                        // 主键ID
	MerchantId   string      `json:"merchantId"   orm:"merchant_id"   description:"商户ID"`                                        // 商户ID
	UserId       string      `json:"userId"       orm:"user_id"       description:"用户ID"`                                        // 用户ID
	PeriodNumber string      `json:"periodNumber" orm:"period_number" description:"期号"`                                          // 期号
	BetType      string      `json:"betType"      orm:"bet_type"      description:"投注类型：high/low/odd/even/num_0等"`               // 投注类型：high/low/odd/even/num_0等
	BetName      string      `json:"betName"      orm:"bet_name"      description:"投注名称：High/Low/Number 0等"`                     // 投注名称：High/Low/Number 0等
	Amount       float64     `json:"amount"       orm:"amount"        description:"投注金额"`                                        // 投注金额
	Multiplier   float64     `json:"multiplier"   orm:"multiplier"    description:"投注时的赔率"`                                      // 投注时的赔率
	Status       string      `json:"status"       orm:"status"        description:"状态：pending-等待开奖，win-已中奖，lose-未中奖，cancel-已取消"` // 状态：pending-等待开奖，win-已中奖，lose-未中奖，cancel-已取消
	Ip           string      `json:"ip"           orm:"ip"            description:"投注IP地址"`                                      // 投注IP地址
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`                                        // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"更新时间"`                                        // 更新时间
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:"删除时间"`                                        // 删除时间
}
