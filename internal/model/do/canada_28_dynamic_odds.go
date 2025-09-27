// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Canada28DynamicOdds is the golang structure of table game_canada28_dynamic_odds for DAO operations like Where/Data.
type Canada28DynamicOdds struct {
	g.Meta             `orm:"table:game_canada28_dynamic_odds, do:true"`
	Id                 interface{} //
	MerchantId         interface{} // 商户ID
	RuleName           interface{} // 规则名称
	TriggerCondition   interface{} // 触发条件：sum_range, sum_exact, sum_in
	TriggerValues      interface{} // 触发条件值（JSON格式）
	BetTypeAdjustments interface{} // 投注类型赔率调整（JSON格式）
	Status             interface{} // 状态：1启用，0禁用
	Priority           interface{} // 优先级
	CreatedAt          *gtime.Time //
	UpdatedAt          *gtime.Time //
}
