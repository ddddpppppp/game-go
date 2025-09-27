// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserBalancesDao is the data access object for the table game_user_balances.
type UserBalancesDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  UserBalancesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// UserBalancesColumns defines and stores column names for the table game_user_balances.
type UserBalancesColumns struct {
	Id            string // ID
	UserId        string // 用户ID
	Type          string // 变动类型：deposit-充值, withdraw-提现, game_bet-投注, game_win-收益
	Amount        string // 变动金额
	BalanceBefore string // 变动前余额
	BalanceAfter  string // 变动后余额
	Description   string // 描述
	RelatedId     string // 关联ID
	CreatedAt     string // 创建时间
}

// userBalancesColumns holds the columns for the table game_user_balances.
var userBalancesColumns = UserBalancesColumns{
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

// NewUserBalancesDao creates and returns a new DAO object for table data access.
func NewUserBalancesDao(handlers ...gdb.ModelHandler) *UserBalancesDao {
	return &UserBalancesDao{
		group:    "default",
		table:    "game_user_balances",
		columns:  userBalancesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserBalancesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserBalancesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserBalancesDao) Columns() UserBalancesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserBalancesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserBalancesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserBalancesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
