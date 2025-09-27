// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GroupMessageDao is the data access object for the table game_group_message.
type GroupMessageDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  GroupMessageColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// GroupMessageColumns defines and stores column names for the table game_group_message.
type GroupMessageColumns struct {
	Id        string // 主键ID
	UserId    string // 发送者ID
	GroupId   string // 群组ID
	Message   string // 消息内容
	Type      string // 消息类型：text-文本，image-图片
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
	DeletedAt string // 删除时间
}

// groupMessageColumns holds the columns for the table game_group_message.
var groupMessageColumns = GroupMessageColumns{
	Id:        "id",
	UserId:    "user_id",
	GroupId:   "group_id",
	Message:   "message",
	Type:      "type",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewGroupMessageDao creates and returns a new DAO object for table data access.
func NewGroupMessageDao(handlers ...gdb.ModelHandler) *GroupMessageDao {
	return &GroupMessageDao{
		group:    "default",
		table:    "game_group_message",
		columns:  groupMessageColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GroupMessageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GroupMessageDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GroupMessageDao) Columns() GroupMessageColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GroupMessageDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GroupMessageDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *GroupMessageDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
