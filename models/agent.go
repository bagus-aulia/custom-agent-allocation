package models

import (
	"github.com/jinzhu/gorm"
)

//Agent models
type Agent struct {
	gorm.Model
	Rooms    []Room `gorm:"foreignkey:AgentID"`
	Nama     string
	Email    string `gorm:"unique_index"`
	SocialID string
	Provider string
	Avatar   string
}
