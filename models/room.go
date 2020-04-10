package models

import (
	"github.com/jinzhu/gorm"
)

//Room models
type Room struct {
	gorm.Model
	Messages   []Message `gorm:"foreignkey:RoomID"`
	Agent      Agent     `gorm:"foreignkey:ID;association_foreignkey:AgentID"`
	Customer   Customer  `gorm:"foreignkey:ID;association_foreignkey:CustomerID"`
	AgentID    uint      `sql:"type:integer REFERENCES rooms(id)"`
	CustomerID uint      `sql:"type:integer REFERENCES members(id)"`
	IsResolved bool      `gorm:"default:0"`
}
