package routes

import (
	"fmt"

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
	// var message models.Message
	rooms := []models.Room{}
	agentID := uint(c.MustGet("jwt_user_id").(float64))

	config.DB.Preload("Messages", "is_readed = ? AND sender_agent = ?", false, false).Find(&rooms, "agent_id", agentID)

}

//AgentLogin to generate agent token
func AgentLogin(c *gin.Context) {
	var agent models.Agent
	name := c.PostForm("name")
	email := c.PostForm("email")

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
	}

	jwtToken := createTokenAgent(&agent)

	c.JSON(200, gin.H{
		"data":    agent,
		"token":   jwtToken,
		"message": "login success",
	})
}
