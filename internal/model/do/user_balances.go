// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserBalances is the golang structure of table game_user_balances for DAO operations like Where/Data.
type UserBalances struct {
	g.Meta        `orm:"table:game_user_balances, do:true"`
	Id            interface{} // ID
	UserId        interface{} // 用户ID
	Type          interface{} // 变动类型：deposit-充值, withdraw-提现, game_bet-投注, game_win-收益
	Amount        interface{} // 变动金额
	BalanceBefore interface{} // 变动前余额
	BalanceAfter  interface{} // 变动后余额
	Description   interface{} // 描述
	RelatedId     interface{} // 关联ID
	CreatedAt     *gtime.Time // 创建时间
}
