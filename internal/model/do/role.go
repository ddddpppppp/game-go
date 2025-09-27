// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure of table game_role for DAO operations like Where/Data.
type Role struct {
	g.Meta     `orm:"table:game_role, do:true"`
	Id         interface{} //
	Name       interface{} //
	Type       interface{} // 1:管理员,2:商户,3:代理,4:个码管理者
	MerchantId interface{} // 商户id
	Access     interface{} //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
}
