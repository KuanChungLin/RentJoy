package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderDetail struct {
	ID        uint            `gorm:"primaryKey;column:Id"`
	OrderID   uint            `gorm:"column:OrderId"`
	StartTime time.Time       `gorm:"column:StartTime"`
	EndTime   time.Time       `gorm:"column:EndTime"`
	Price     decimal.Decimal `gorm:"column:Price"`

	Order Order `gorm:"foreignKey:OrderID"`
}

func (OrderDetail) TableName() string {
	return "OrderDetail"
}
