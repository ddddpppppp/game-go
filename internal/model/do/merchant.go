// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Merchant is the golang structure of table game_merchant for DAO operations like Where/Data.
type Merchant struct {
	g.Meta    `orm:"table:game_merchant, do:true"`
	Id        interface{} //
	AdminId   interface{} // admin表的id
	Balance   interface{} // 余额
	Uuid      interface{} //
	Name      interface{} //
	Logo      interface{} //
	AppKey    interface{} //
	Type      interface{} //
	Status    interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}
