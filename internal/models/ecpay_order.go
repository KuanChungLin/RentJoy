package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type EcpayOrder struct {
	ID               uint            `gorm:"primaryKey;column:Id"`
	MerchantID       string          `gorm:"column:MerchantId"`
	MerchantTradeNum string          `gorm:"column:MerchantTradeNO"`
	RtnCode          int             `gorm:"column:RtnCode"`
	RtnMsg           string          `gorm:"column:RtnMsg"`
	TradeNum         int             `gorm:"column:TradeNO"`
	TradeAmt         decimal.Decimal `gorm:"column:TradeAmt"`
	PaymentDate      time.Time       `gorm:"column:PaymentDate"`
	PaymentType      string          `gorm:"column:PaymentType"`
	Charge           decimal.Decimal `gorm:"column:PaymentTypeChargeFree"`
	TradeDate        time.Time       `gorm:"column:TradeDate"`
	SimulatePaid     bool            `gorm:"column:SimulatePaid"`
	CheckCode        string          `gorm:"column:CheckMacValue"`

	Order *Order `gorm:"foreignKey:ID"`
}

func (EcpayOrder) TableName() string {
	return "EcpayOrder"
}
