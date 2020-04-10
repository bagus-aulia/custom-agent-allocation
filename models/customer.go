package models

import (
	"github.com/jinzhu/gorm"
)

//Customer models
type Customer struct {
	gorm.Model
	Rooms  []Room `gorm:"foreignkey:CustomerID"`
	Name   string
	Email  string `gorm:"unique_index"`
	Avatar string
}
