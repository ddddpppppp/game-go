// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupMessage is the golang structure of table game_group_message for DAO operations like Where/Data.
type GroupMessage struct {
	g.Meta    `orm:"table:game_group_message, do:true"`
	Id        interface{} // 主键ID
	UserId    interface{} // 发送者ID
	GroupId   interface{} // 群组ID
	Message   interface{} // 消息内容
	Type      interface{} // 消息类型：text-文本，image-图片
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
	DeletedAt *gtime.Time // 删除时间
}
