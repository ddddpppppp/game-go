// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Bingo28BetTypesDao is the data access object for the table game_bingo28_bet_types.
type Bingo28BetTypesDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  Bingo28BetTypesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// Bingo28BetTypesColumns defines and stores column names for the table game_bingo28_bet_types.
type Bingo28BetTypesColumns struct {
	Id          string //
	MerchantId  string // 商户ID
	TypeName    string // 玩法名称
	TypeKey     string // 玩法标识
	Description string // 玩法描述
	Odds        string // 赔率倍数
	Status      string // 状态：1启用，0禁用
	Sort        string // 排序
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// bingo28BetTypesColumns holds the columns for the table game_bingo28_bet_types.
var bingo28BetTypesColumns = Bingo28BetTypesColumns{
	Id:          "id",
	MerchantId:  "merchant_id",
	TypeName:    "type_name",
	TypeKey:     "type_key",
	Description: "description",
	Odds:        "odds",
	Status:      "status",
	Sort:        "sort",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewBingo28BetTypesDao creates and returns a new DAO object for table data access.
func NewBingo28BetTypesDao(handlers ...gdb.ModelHandler) *Bingo28BetTypesDao {
	return &Bingo28BetTypesDao{
		group:    "default",
		table:    "game_bingo28_bet_types",
		columns:  bingo28BetTypesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *Bingo28BetTypesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *Bingo28BetTypesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *Bingo28BetTypesDao) Columns() Bingo28BetTypesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *Bingo28BetTypesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *Bingo28BetTypesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *Bingo28BetTypesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
