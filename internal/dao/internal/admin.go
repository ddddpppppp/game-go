// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AdminDao is the data access object for the table game_admin.
type AdminDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  AdminColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// AdminColumns defines and stores column names for the table game_admin.
type AdminColumns struct {
	Id         string //
	Uuid       string //
	Nickname   string //
	Avatar     string //
	Username   string //
	Balance    string // 余额
	Password   string //
	Salt       string //
	MerchantId string // 对应merchant表
	RoleId     string // 对应role表
	ParentId   string // 推荐人ID
	Path       string // 路径（记录所有上级ID包括自己，如0:1:2:3）
	Depth      string // 层级深度
	Status     string // -1冻结，1开启
	CreatedAt  string //
	UpdatedAt  string //
	DeletedAt  string //
}

// adminColumns holds the columns for the table game_admin.
var adminColumns = AdminColumns{
	Id:         "id",
	Uuid:       "uuid",
	Nickname:   "nickname",
	Avatar:     "avatar",
	Username:   "username",
	Balance:    "balance",
	Password:   "password",
	Salt:       "salt",
	MerchantId: "merchant_id",
	RoleId:     "role_id",
	ParentId:   "parent_id",
	Path:       "path",
	Depth:      "depth",
	Status:     "status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewAdminDao creates and returns a new DAO object for table data access.
func NewAdminDao(handlers ...gdb.ModelHandler) *AdminDao {
	return &AdminDao{
		group:    "default",
		table:    "game_admin",
		columns:  adminColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AdminDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AdminDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AdminDao) Columns() AdminColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AdminDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AdminDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AdminDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
