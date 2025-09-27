// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentChannelDao is the data access object for the table game_payment_channel.
type PaymentChannelDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  PaymentChannelColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// PaymentChannelColumns defines and stores column names for the table game_payment_channel.
type PaymentChannelColumns struct {
	Id            string //
	Name          string // 名称
	BelongAdminId string // 所属管理者id
	Type          string // 类型
	Rate          string // 费率
	ChargeFee     string // 单笔手续费
	CountTime     string // 结算时间
	Guarantee     string // 保证金
	FreezeTime    string // 冻结时间
	DayLimitMoney string // 每日限额
	DayLimitCount string // 每日限额次数
	Remark        string // 备注
	Status        string // -1停用，1开启
	IsBackup      string // 是否备用渠道
	Params        string // 渠道参数
	Sort          string //
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
}

// paymentChannelColumns holds the columns for the table game_payment_channel.
var paymentChannelColumns = PaymentChannelColumns{
	Id:            "id",
	Name:          "name",
	BelongAdminId: "belong_admin_id",
	Type:          "type",
	Rate:          "rate",
	ChargeFee:     "charge_fee",
	CountTime:     "count_time",
	Guarantee:     "guarantee",
	FreezeTime:    "freeze_time",
	DayLimitMoney: "day_limit_money",
	DayLimitCount: "day_limit_count",
	Remark:        "remark",
	Status:        "status",
	IsBackup:      "is_backup",
	Params:        "params",
	Sort:          "sort",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewPaymentChannelDao creates and returns a new DAO object for table data access.
func NewPaymentChannelDao(handlers ...gdb.ModelHandler) *PaymentChannelDao {
	return &PaymentChannelDao{
		group:    "default",
		table:    "game_payment_channel",
		columns:  paymentChannelColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PaymentChannelDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PaymentChannelDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PaymentChannelDao) Columns() PaymentChannelColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PaymentChannelDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PaymentChannelDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PaymentChannelDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
