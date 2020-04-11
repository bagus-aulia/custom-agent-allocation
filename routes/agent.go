package routes

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/bagus-aulia/custom-agent-allocation/config"
	"github.com/bagus-aulia/custom-agent-allocation/models"
	"github.com/gin-gonic/gin"
)

//AgentIndex to view chat list
func AgentIndex(c *gin.Context) {
	rooms := []models.Room{}
	agentID := uint(c.MustGet("jwt_user_id").(float64))

	fmt.Println(agentID)

	config.DB.Preload("Messages").Find(&rooms, "agent_id = ?", agentID)

	c.JSON(200, gin.H{
		"data": rooms,
	})
}

//AgentReceive to get unread message from customer
func AgentReceive(c *gin.Context) {
	rooms := []models.Room{}
	agentID := uint(c.MustGet("jwt_user_id").(float64))

	config.DB.Preload("Messages", "is_readed = ? AND sender_agent = ?", false, false).Find(&rooms, "agent_id = ?", agentID)

	c.JSON(200, gin.H{
		"data": rooms,
	})
}

//AgentLogin to generate agent token
func AgentLogin(c *gin.Context) {
	var agent models.Agent
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := getMD5Hash(c.PostForm("password"))

	if email == "" {
		c.JSON(403, gin.H{
			"message": "email must be filled",
		})

		c.Abort()
		return
	}

	if config.DB.First(&agent, "email = ?", email).RecordNotFound() {
		agent = models.Agent{
			Name:   name,
			Email:  email,
			Avatar: "default.png",
		}

		config.DB.Create(&agent)
	} else {
		if agent.Password != password {
			c.JSON(403, gin.H{
				"message": "password salah",
				"name":    name,
				"email":   email,
			})

			c.Abort()
			return
		}
	}

	jwtToken := createTokenAgent(&agent)

	c.JSON(200, gin.H{
		"data":    agent,
		"token":   jwtToken,
		"message": "login success",
	})
}

//AgentSend to save chat
func AgentSend(c *gin.Context) {
	var room models.Room
	agentID := int(c.MustGet("jwt_user_id").(float64))
	roomID := c.PostForm("room_id")

	if config.DB.First(&room, roomID).RecordNotFound() {
		c.JSON(404, gin.H{
			"message": "room not available",
		})

		c.Abort()
		return
	}

	roomIDuint, _ := strconv.ParseUint(c.PostForm("room_id"), 10, 32)
	message := models.Message{
		RoomID:      uint(roomIDuint),
		SenderID:    agentID,
		SenderAgent: true,
		Message:     c.PostForm("message"),
	}

	config.DB.Create(&message)

	c.JSON(200, gin.H{
		"data": message,
	})
}

//AgentRead to set readed message status
func AgentRead(c *gin.Context) {
	var room models.Room
	message := []models.Message{}
	agentID := uint(c.MustGet("jwt_user_id").(float64))
	roomID := c.PostForm("room_id")
	returnMsg := ""

	if config.DB.First(&room, roomID).RecordNotFound() {
		returnMsg = "room doesn't exist"
	} else {
		if room.AgentID != agentID {
			returnMsg = "you haven't access to this room"
		}
	}

	if returnMsg != "" {
		c.JSON(404, gin.H{
			"message": returnMsg,
		})

		c.Abort()
		return
	}

	config.DB.Model(&message).Where("room_id = ? AND sender_agent = ? AND is_readed = ?", room.ID, false, false).Update("is_readed", true)

	c.JSON(200, gin.H{
		"message": "message readed",
	})
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
