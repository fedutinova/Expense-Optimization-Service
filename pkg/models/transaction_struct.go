package models

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

var (
	incomeCategory = map[string]struct{}{
		"salary": {},
		"gift":   {},
		"bonus":  {},
		"loan":   {},
		"other":  {},
	}
	expensesCats = map[string]struct{}{
		"food":               {},
		"transportation":     {},
		"housing":            {},
		"entertainment":      {},
		"fashion":            {},
		"tourism":            {},
		"education":          {},
		"insurance":          {},
		"beauty and health":  {},
		"animals":            {},
		"services":           {},
		"loans and payments": {},
		"other":              {},
	}
)

type Transaction struct {
	ID          int             `json:"-" db:"id"`
	UserId      int             `json:"user_id" db:"user_id"`
	Date        time.Time       `json:"date" db:"date"`
	Description string          `json:"description" db:"description"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	Category    string          `json:"category" db:"category"`
}

func (i Transaction) Validate() error {
	if i.Amount.GreaterThan(decimal.Zero) {
		if _, found := incomeCategory[i.Category]; found {
			return nil
		}

		return errors.New("invalid category")
	} else {
		if _, found := expensesCats[i.Category]; found {
			return nil
		}

		return errors.New("invalid category")
	}
}

type TransactionsByCategory struct {
	Category string          `json:"category"`
	Sum      decimal.Decimal `json:"sum" db:"sum"`
}

type TransactionDate struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type UpdateTransactionInput struct {
	Date        *time.Time      `json:"date"`
	Description *string         `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	Category    *string         `json:"category"`
}

func (i UpdateTransactionInput) Validate() error {
	if i.Date == nil && i.Description == nil && i.Amount == (decimal.Decimal{}) && i.Category == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
