// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TransactionsDao is the data access object for the table game_transactions.
type TransactionsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  TransactionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// TransactionsColumns defines and stores column names for the table game_transactions.
type TransactionsColumns struct {
	Id           string // ID
	UserId       string // 用户ID
	Type         string // 交易类型: deposit-充值, withdraw-提现
	ChannelId    string // 支付渠道ID
	Amount       string // 交易金额
	ActualAmount string // 实际金额
	Account      string // 账户
	OrderNo      string // 订单号
	Fee          string // 手续费
	Gift         string // 赠送金额
	Status       string // 状态: pending-待处理, completed-已完成, failed-失败, expired-已过期
	CreatedAt    string // 创建时间
	CompletedAt  string // 完成时间
	ExpiredAt    string // 过期时间
	DeletedAt    string // 删除时间
}

// transactionsColumns holds the columns for the table game_transactions.
var transactionsColumns = TransactionsColumns{
	Id:           "id",
	UserId:       "user_id",
	Type:         "type",
	ChannelId:    "channel_id",
	Amount:       "amount",
	ActualAmount: "actual_amount",
	Account:      "account",
	OrderNo:      "order_no",
	Fee:          "fee",
	Gift:         "gift",
	Status:       "status",
	CreatedAt:    "created_at",
	CompletedAt:  "completed_at",
	ExpiredAt:    "expired_at",
	DeletedAt:    "deleted_at",
}

// NewTransactionsDao creates and returns a new DAO object for table data access.
func NewTransactionsDao(handlers ...gdb.ModelHandler) *TransactionsDao {
	return &TransactionsDao{
		group:    "default",
		table:    "game_transactions",
		columns:  transactionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TransactionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TransactionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TransactionsDao) Columns() TransactionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TransactionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TransactionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TransactionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
