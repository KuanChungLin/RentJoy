package models

import (
	"github.com/shopspring/decimal"
)

type Budget struct {
	ID    uint            `gorm:"primaryKey;column:Id"`
	Price decimal.Decimal `gorm:"column:Price"`
}

func (Budget) TableName() string {
	return "Budget"
}
