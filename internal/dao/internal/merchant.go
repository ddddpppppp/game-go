// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantDao is the data access object for the table game_merchant.
type MerchantDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  MerchantColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// MerchantColumns defines and stores column names for the table game_merchant.
type MerchantColumns struct {
	Id        string //
	AdminId   string // admin表的id
	Balance   string // 余额
	Uuid      string //
	Name      string //
	Logo      string //
	AppKey    string //
	Type      string //
	Status    string //
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
}

// merchantColumns holds the columns for the table game_merchant.
var merchantColumns = MerchantColumns{
	Id:        "id",
	AdminId:   "admin_id",
	Balance:   "balance",
	Uuid:      "uuid",
	Name:      "name",
	Logo:      "logo",
	AppKey:    "app_key",
	Type:      "type",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewMerchantDao creates and returns a new DAO object for table data access.
func NewMerchantDao(handlers ...gdb.ModelHandler) *MerchantDao {
	return &MerchantDao{
		group:    "default",
		table:    "game_merchant",
		columns:  merchantColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MerchantDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MerchantDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MerchantDao) Columns() MerchantColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MerchantDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MerchantDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MerchantDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
