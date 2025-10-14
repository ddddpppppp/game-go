// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// BingoDynamicOdds is the golang structure for table bingo_dynamic_odds.
type BingoDynamicOdds struct {
	Id                 int         `json:"id"                 orm:"id"                   description:""`                                  //
	MerchantId         string      `json:"merchantId"         orm:"merchant_id"          description:"商户ID"`                              // 商户ID
	RuleName           string      `json:"ruleName"           orm:"rule_name"            description:"规则名称"`                              // 规则名称
	TriggerCondition   string      `json:"triggerCondition"   orm:"trigger_condition"    description:"触发条件：sum_range, sum_exact, sum_in"` // 触发条件：sum_range, sum_exact, sum_in
	TriggerValues      string      `json:"triggerValues"      orm:"trigger_values"       description:"触发条件值（JSON格式）"`                     // 触发条件值（JSON格式）
	BetTypeAdjustments string      `json:"betTypeAdjustments" orm:"bet_type_adjustments" description:"投注类型赔率调整（JSON格式）"`                  // 投注类型赔率调整（JSON格式）
	Status             int         `json:"status"             orm:"status"               description:"状态：1启用，0禁用"`                        // 状态：1启用，0禁用
	Priority           int         `json:"priority"           orm:"priority"             description:"优先级"`                               // 优先级
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           description:""`                                  //
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           description:""`                                  //
}
