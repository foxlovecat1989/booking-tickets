package db

import "github.com/shopspring/decimal"

type Order struct {
	ID         int             `db:"id"`
	CreatedAt  int64           `db:"created_at"`
	Status     string          `db:"status"`
	TotalPrice decimal.Decimal `db:"total_price"`
}
