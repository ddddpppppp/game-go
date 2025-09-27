// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Canada28BetTypes is the golang structure of table game_canada28_bet_types for DAO operations like Where/Data.
type Canada28BetTypes struct {
	g.Meta      `orm:"table:game_canada28_bet_types, do:true"`
	Id          interface{} //
	MerchantId  interface{} // 商户ID
	TypeName    interface{} // 玩法名称
	TypeKey     interface{} // 玩法标识
	Description interface{} // 玩法描述
	Odds        interface{} // 赔率倍数
	Status      interface{} // 状态：1启用，0禁用
	Sort        interface{} // 排序
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
