package models

import (
	"github.com/jinzhu/gorm"
)

//Message models
type Message struct {
	gorm.Model
	Room     Room `gorm:"foreignkey:ID;association_foreignkey:RoomID"`
	RoomID   uint `sql:"type:integer REFERENCES room(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	SenderID uint
	Message  string `sql:"type:text;"`
}
