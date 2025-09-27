// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupMessage is the golang structure for table group_message.
type GroupMessage struct {
	Id        uint64      `json:"id"        orm:"id"         description:"主键ID"`                  // 主键ID
	UserId    string      `json:"userId"    orm:"user_id"    description:"发送者ID"`                 // 发送者ID
	GroupId   string      `json:"groupId"   orm:"group_id"   description:"群组ID"`                  // 群组ID
	Message   string      `json:"message"   orm:"message"    description:"消息内容"`                  // 消息内容
	Type      string      `json:"type"      orm:"type"       description:"消息类型：text-文本，image-图片"` // 消息类型：text-文本，image-图片
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`                  // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`                  // 更新时间
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"删除时间"`                  // 删除时间
}
