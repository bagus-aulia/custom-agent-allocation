package models

import (
	"github.com/jinzhu/gorm"
)

//Message models
type Message struct {
	gorm.Model
	RoomID      uint
	SenderID    int
	SenderAgent bool   `gorm:"default:0"`
	Message     string `sql:"type:text;"`
	IsReaded    bool   `gorm:"default:0"`
}
