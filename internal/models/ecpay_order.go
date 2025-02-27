package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type EcpayOrder struct {
	ID              uint            `gorm:"primaryKey;column:Id;autoIncrement:false"` // 取消自動增長以因應反關聯性資料
	MerchantID      string          `gorm:"column:MerchantId"`
	MerchantTradeNo string          `gorm:"column:MerchantTradeNo"`
	RtnCode         int             `gorm:"column:RtnCode"`
	RtnMsg          string          `gorm:"column:RtnMsg"`
	TradeNo         int             `gorm:"column:TradeNo"`
	TradeAmt        decimal.Decimal `gorm:"column:TradeAmt"`
	PaymentDate     time.Time       `gorm:"column:PaymentDate"`
	PaymentType     string          `gorm:"column:PaymentType"`
	Charge          decimal.Decimal `gorm:"column:PaymentTypeChargeFree"`
	TradeDate       time.Time       `gorm:"column:TradeDate"`
	SimulatePaid    bool            `gorm:"column:SimulatePaid"`
	CheckMacValue   string          `gorm:"column:CheckMacValue"`

	Order *Order `gorm:"foreignKey:ID"`
}

func (EcpayOrder) TableName() string {
	return "EcpayOrder"
}
