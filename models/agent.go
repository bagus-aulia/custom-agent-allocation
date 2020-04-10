package models

import (
	"github.com/jinzhu/gorm"
)

//Agent models
type Agent struct {
	gorm.Model
	Rooms    []Room `gorm:"foreignkey:AgentID"`
	Name     string
	Email    string `gorm:"unique_index"`
	Password string
	Avatar   string
}
