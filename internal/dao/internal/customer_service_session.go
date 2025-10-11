// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerServiceSessionDao is the data access object for the table game_customer_service_session.
type CustomerServiceSessionDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  CustomerServiceSessionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// CustomerServiceSessionColumns defines and stores column names for the table game_customer_service_session.
type CustomerServiceSessionColumns struct {
	Id            string // 主键ID
	UserId        string // 用户ID
	AdminId       string // 当前服务的管理员ID
	LastMessage   string // 最后一条消息
	LastMessageAt string // 最后消息时间
	UnreadCount   string // 未读消息数
	Status        string // 会话状态：1-活跃，2-已关闭
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
	DeletedAt     string // 删除时间
}

// customerServiceSessionColumns holds the columns for the table game_customer_service_session.
var customerServiceSessionColumns = CustomerServiceSessionColumns{
	Id:            "id",
	UserId:        "user_id",
	AdminId:       "admin_id",
	LastMessage:   "last_message",
	LastMessageAt: "last_message_at",
	UnreadCount:   "unread_count",
	Status:        "status",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewCustomerServiceSessionDao creates and returns a new DAO object for table data access.
func NewCustomerServiceSessionDao(handlers ...gdb.ModelHandler) *CustomerServiceSessionDao {
	return &CustomerServiceSessionDao{
		group:    "default",
		table:    "game_customer_service_session",
		columns:  customerServiceSessionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CustomerServiceSessionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CustomerServiceSessionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CustomerServiceSessionDao) Columns() CustomerServiceSessionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CustomerServiceSessionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CustomerServiceSessionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CustomerServiceSessionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
