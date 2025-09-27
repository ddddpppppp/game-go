// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersDao is the data access object for the table game_users.
type UsersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UsersColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UsersColumns defines and stores column names for the table game_users.
type UsersColumns struct {
	Id            string //
	Uuid          string // 用户名/手机号/邮箱
	Username      string // 用户名/手机号/邮箱
	Type          string // bot/user
	Password      string // 加密后的密码
	Balance       string // 余额
	BalanceFrozen string // 冻结余额
	Nickname      string // 昵称
	Avatar        string // 头像URL
	MerchantId    string // 商户ID
	ParentId      string // 邀请人ID
	Status        string // 状态 (1:正常, 0:禁用)
	Salt          string // 密码盐值
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
}

// usersColumns holds the columns for the table game_users.
var usersColumns = UsersColumns{
	Id:            "id",
	Uuid:          "uuid",
	Username:      "username",
	Type:          "type",
	Password:      "password",
	Balance:       "balance",
	BalanceFrozen: "balance_frozen",
	Nickname:      "nickname",
	Avatar:        "avatar",
	MerchantId:    "merchant_id",
	ParentId:      "parent_id",
	Status:        "status",
	Salt:          "salt",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewUsersDao creates and returns a new DAO object for table data access.
func NewUsersDao(handlers ...gdb.ModelHandler) *UsersDao {
	return &UsersDao{
		group:    "default",
		table:    "game_users",
		columns:  usersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UsersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UsersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UsersDao) Columns() UsersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UsersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UsersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UsersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
