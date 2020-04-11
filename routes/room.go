package routes

import (
	"github.com/bagus-aulia/custom-agent-allocation/config"
	"github.com/bagus-aulia/custom-agent-allocation/models"
	"github.com/gin-gonic/gin"
)

//RoomDivide to divide unhandled message to available agent
func RoomDivide(c *gin.Context) {
	agents := []models.Agent{}
	rooms := []models.Room{}
	var antiranAgent models.Agent
	var agent models.Agent
	var room models.Room

	config.DB.First(&antiranAgent, "name = ?", "antrian")

	for i := 0; i < 2; i++ {
		config.DB.Where("name <> ?", "antrian").Find(&agents, "handled BETWEEN ? AND ?", 0, 1)
		config.DB.Order("updated_at asc").Where("agent_id = ?", antiranAgent.ID).Find(&rooms)

		for i, roomData := range rooms {
			currentAgent := agents[i]

			config.DB.Model(&room).First(&room, roomData.ID).Update("agent_id", agent.ID)
			config.DB.Model(&agent).First(&agent, currentAgent.ID).Update("handled", (currentAgent.Handled + 1))

			if i == (len(agents) - 1) {
				break
			}
		}
	}

	c.JSON(200, gin.H{
		"message": "room has been divided",
	})
}

//RoomResolve to set room as resolved
func RoomResolve(c *gin.Context) {
	var agent models.Agent
	var room models.Room
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

	config.DB.First(&agent, agentID)
	config.DB.Model(&room).First(&room, roomID).Update("is_resolved", true)
	config.DB.Model(&agent).First(&agent, agentID).Update("handled", (agent.Handled - 1))

	c.JSON(200, gin.H{
		"message": "room " + roomID + " has been resolved",
	})
}
