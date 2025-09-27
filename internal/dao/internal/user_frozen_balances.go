// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserFrozenBalancesDao is the data access object for the table game_user_frozen_balances.
type UserFrozenBalancesDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  UserFrozenBalancesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// UserFrozenBalancesColumns defines and stores column names for the table game_user_frozen_balances.
type UserFrozenBalancesColumns struct {
	Id            string // ID
	UserId        string // 用户ID
	Type          string // 变动类型：game_bet-投注
	Amount        string // 变动金额
	BalanceBefore string // 变动前余额
	BalanceAfter  string // 变动后余额
	Description   string // 描述
	RelatedId     string // 关联ID
	CreatedAt     string // 创建时间
}

// userFrozenBalancesColumns holds the columns for the table game_user_frozen_balances.
var userFrozenBalancesColumns = UserFrozenBalancesColumns{
	Id:            "id",
	UserId:        "user_id",
	Type:          "type",
	Amount:        "amount",
	BalanceBefore: "balance_before",
	BalanceAfter:  "balance_after",
	Description:   "description",
	RelatedId:     "related_id",
	CreatedAt:     "created_at",
}

// NewUserFrozenBalancesDao creates and returns a new DAO object for table data access.
func NewUserFrozenBalancesDao(handlers ...gdb.ModelHandler) *UserFrozenBalancesDao {
	return &UserFrozenBalancesDao{
		group:    "default",
		table:    "game_user_frozen_balances",
		columns:  userFrozenBalancesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserFrozenBalancesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserFrozenBalancesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserFrozenBalancesDao) Columns() UserFrozenBalancesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserFrozenBalancesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserFrozenBalancesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserFrozenBalancesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
