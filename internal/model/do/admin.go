// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Admin is the golang structure of table game_admin for DAO operations like Where/Data.
type Admin struct {
	g.Meta     `orm:"table:game_admin, do:true"`
	Id         interface{} //
	Uuid       interface{} //
	Nickname   interface{} //
	Avatar     interface{} //
	Username   interface{} //
	Balance    interface{} // 余额
	Password   interface{} //
	Salt       interface{} //
	MerchantId interface{} // 对应merchant表
	RoleId     interface{} // 对应role表
	ParentId   interface{} // 推荐人ID
	Path       interface{} // 路径（记录所有上级ID包括自己，如0:1:2:3）
	Depth      interface{} // 层级深度
	Status     interface{} // -1冻结，1开启
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
}
