// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserFrozenBalances is the golang structure for table user_frozen_balances.
type UserFrozenBalances struct {
	Id            uint64      `json:"id"            orm:"id"             description:"ID"`               // ID
	UserId        uint64      `json:"userId"        orm:"user_id"        description:"用户ID"`             // 用户ID
	Type          string      `json:"type"          orm:"type"           description:"变动类型：game_bet-投注"` // 变动类型：game_bet-投注
	Amount        float64     `json:"amount"        orm:"amount"         description:"变动金额"`             // 变动金额
	BalanceBefore float64     `json:"balanceBefore" orm:"balance_before" description:"变动前余额"`            // 变动前余额
	BalanceAfter  float64     `json:"balanceAfter"  orm:"balance_after"  description:"变动后余额"`            // 变动后余额
	Description   string      `json:"description"   orm:"description"    description:"描述"`               // 描述
	RelatedId     string      `json:"relatedId"     orm:"related_id"     description:"关联ID"`             // 关联ID
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:"创建时间"`             // 创建时间
}
