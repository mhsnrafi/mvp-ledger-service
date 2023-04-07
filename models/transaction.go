package models

import (
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	UserID        uint
	Amount        float64
	Type          string // "credit" or "debit"
	TransactionID string `gorm:"unique;not null"`
}
