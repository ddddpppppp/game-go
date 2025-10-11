// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KenoBets is the golang structure for table keno_bets.
type KenoBets struct {
	Id              uint64      `json:"id"              orm:"id"               description:"主键ID"`                                        // 主键ID
	MerchantId      string      `json:"merchantId"      orm:"merchant_id"      description:"商户ID"`                                        // 商户ID
	UserId          string      `json:"userId"          orm:"user_id"          description:"用户ID"`                                        // 用户ID
	PeriodNumber    string      `json:"periodNumber"    orm:"period_number"    description:"期号"`                                          // 期号
	SelectedNumbers string      `json:"selectedNumbers" orm:"selected_numbers" description:"玩家选择的10个号码，JSON格式: [1,5,12,...]"`             // 玩家选择的10个号码，JSON格式: [1,5,12,...]
	DrawnNumbers    string      `json:"drawnNumbers"    orm:"drawn_numbers"    description:"开出的20个号码，JSON格式: [3,7,12,...]"`               // 开出的20个号码，JSON格式: [3,7,12,...]
	MatchedNumbers  string      `json:"matchedNumbers"  orm:"matched_numbers"  description:"匹配的号码，JSON格式: [12,...]"`                      // 匹配的号码，JSON格式: [12,...]
	MatchCount      int         `json:"matchCount"      orm:"match_count"      description:"匹配数量：0-10"`                                   // 匹配数量：0-10
	Amount          float64     `json:"amount"          orm:"amount"           description:"投注金额"`                                        // 投注金额
	Multiplier      float64     `json:"multiplier"      orm:"multiplier"       description:"投注时的赔率（基于匹配数量）"`                              // 投注时的赔率（基于匹配数量）
	WinAmount       float64     `json:"winAmount"       orm:"win_amount"       description:"中奖金额"`                                        // 中奖金额
	Status          string      `json:"status"          orm:"status"           description:"状态：pending-等待开奖，win-已中奖，lose-未中奖，cancel-已取消"` // 状态：pending-等待开奖，win-已中奖，lose-未中奖，cancel-已取消
	Ip              string      `json:"ip"              orm:"ip"               description:"投注IP地址"`                                      // 投注IP地址
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:"创建时间"`                                        // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:"更新时间"`                                        // 更新时间
	SettledAt       *gtime.Time `json:"settledAt"       orm:"settled_at"       description:"结算时间"`                                        // 结算时间
	DeletedAt       *gtime.Time `json:"deletedAt"       orm:"deleted_at"       description:"删除时间"`                                        // 删除时间
}
