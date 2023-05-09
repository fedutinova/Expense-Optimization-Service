package models

import (
	"github.com/shopspring/decimal"
)

type Member struct {
	Id       int `json:"user_id"`
	WalletId int `json:"wallet_id"`
}

type Wallet struct {
	ID          int             `json:"-" db:"id"`
	IsFamily    bool            `json:"is_family"`
	Balance     decimal.Decimal `json:"balance"`
	Description string          `json:"description"`
}
