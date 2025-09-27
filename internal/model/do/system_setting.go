// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemSetting is the golang structure of table game_system_setting for DAO operations like Where/Data.
type SystemSetting struct {
	g.Meta      `orm:"table:game_system_setting, do:true"`
	Id          interface{} //
	Name        interface{} // 设置名称，如recharge_setting、withdraw_setting等
	Title       interface{} // 设置标题
	Description interface{} // 设置描述
	Config      interface{} // 设置配置JSON数据
	Status      interface{} // 状态：1启用，0禁用
	Sort        interface{} // 排序
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
