package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UID          string `gorm:"unique;not null"`
	Balance      float64
	Transactions []Transaction
}
