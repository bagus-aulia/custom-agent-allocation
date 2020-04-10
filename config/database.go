package config

import (
	"os"

	"github.com/bagus-aulia/custom-agent-allocation/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //connect mysql
)

//DB is global variable to access database
var DB *gorm.DB

//InitDB used to connect to database function
func InitDB() {
	var err error

	DB, err = gorm.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@/"+os.Getenv("DB_NAME")+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("Failed to connect database")
	}

	DB.AutoMigrate(&models.Agent{})
	DB.AutoMigrate(&models.Customer{})
	DB.AutoMigrate(&models.Room{}).AddForeignKey("customer_id", "customers(id)", "RESTRICT", "RESTRICT").AddForeignKey("agent_id", "agents(id)", "RESTRICT", "RESTRICT")
	DB.AutoMigrate(&models.Message{}).AddForeignKey("room_id", "rooms(id)", "RESTRICT", "RESTRICT")

	DB.Model(&models.Agent{}).Related(&models.Room{})
	DB.Model(&models.Customer{}).Related(&models.Room{})
	DB.Model(&models.Room{}).Related(&models.Message{})
}
