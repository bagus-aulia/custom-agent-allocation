package routes

import (
	"fmt"

	"github.com/bagus-aulia/custom-agent-allocation/config"
	"github.com/bagus-aulia/custom-agent-allocation/models"
	"github.com/gin-gonic/gin"
)

//CustomerIndex to view customer chat
func CustomerIndex(c *gin.Context) {
	var customer models.Customer
	customerID := uint(c.MustGet("jwt_user_id").(float64))

	if config.DB.Preload("Rooms").Preload("Rooms.Messages").Preload("Rooms.Agent").First(&customer, customerID).RecordNotFound() {
		c.JSON(404, gin.H{
			"message": "customer not found",
		})

		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"data": customer,
	})
}

//CustomerReceive to get unread message from agent
func CustomerReceive(c *gin.Context) {
	var room models.Room
	customerID := uint(c.MustGet("jwt_user_id").(float64))

	config.DB.Preload("Messages", "is_readed = ? AND sender_agent = ?", false, true).First(&room, "customer_id = ?", customerID)

	c.JSON(200, gin.H{
		"data": room.Messages,
	})
}

//CustomerLogin to generate customer token
func CustomerLogin(c *gin.Context) {
	var customer models.Customer
	name := c.PostForm("name")
	email := c.PostForm("email")

	if email == "" {
		c.JSON(403, gin.H{
			"message": "email must be filled",
		})

		c.Abort()
		return
	}

	if config.DB.First(&customer, "email = ?", email).RecordNotFound() {
		customer = models.Customer{
			Name:   name,
			Email:  email,
			Avatar: "default.png",
		}

		config.DB.Create(&customer)
	}

	jwtToken := createTokenCustomer(&customer)

	c.JSON(200, gin.H{
		"data":    customer,
		"token":   jwtToken,
		"message": "login success",
	})
}

//CustomerSend to save chat
func CustomerSend(c *gin.Context) {
	var room models.Room
	var antrian models.Agent
	customerID := uint(c.MustGet("jwt_user_id").(float64))

	config.DB.First(&antrian, "name = ?", "antrian")

	if config.DB.First(&room, "customer_id = ?", customerID).RecordNotFound() {
		room = models.Room{
			CustomerID: customerID,
			AgentID:    antrian.ID,
		}

		config.DB.Create(&room)
	} else {
		if room.IsResolved == true {
			if err := config.DB.Model(&room).First(&room, "customer_id = ?", customerID).Updates(map[string]interface{}{
				"AgentID":    antrian.ID,
				"IsResolved": false,
			}).Error; err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	senderID := int(customerID)
	message := models.Message{
		RoomID:   room.ID,
		SenderID: senderID,
		Message:  c.PostForm("message"),
	}

	config.DB.Create(&message)

	c.JSON(200, gin.H{
		"data": message,
	})
}

//CustomerRead to set readed message status
func CustomerRead(c *gin.Context) {
	var room models.Room
	message := []models.Message{}
	customerID := uint(c.MustGet("jwt_user_id").(float64))

	if config.DB.First(&room, "customer_id = ?", customerID).RecordNotFound() {
		c.JSON(404, gin.H{
			"message": "room doesn't exist",
		})

		c.Abort()
		return
	}

	config.DB.Model(&message).Where("room_id = ? AND sender_agent = ? AND is_readed = ?", room.ID, true, false).Update("is_readed", true)

	c.JSON(200, gin.H{
		"message": "message readed",
	})
}
