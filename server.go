package main

import (
	"github.com/bagus-aulia/custom-agent-allocation/config"
	"github.com/bagus-aulia/custom-agent-allocation/middleware"
	"github.com/bagus-aulia/custom-agent-allocation/routes"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	config.InitDB()
	defer config.DB.Close()

	route := gin.Default()

	api := route.Group("/api/v1/")
	{
		customer := api.Group("/customer")
		{
			customer.GET("/", middleware.IsAuth(), routes.CustomerIndex)
			customer.GET("/receive", middleware.IsAuth(), routes.CustomerReceive)
			customer.GET("/readed", middleware.IsAuth(), routes.CustomerRead)

			customer.POST("/login", routes.CustomerLogin)
			customer.POST("/send", middleware.IsAuth(), routes.CustomerSend)
		}

		agent := api.Group("agent")
		{
			agent.GET("/", middleware.IsAdmin(), routes.AgentIndex)
			agent.GET("/receive", middleware.IsAdmin(), routes.AgentReceive)

			agent.POST("/login", routes.AgentLogin)
			agent.POST("/send", middleware.IsAdmin(), routes.AgentSend)

			agent.PUT("/readed", middleware.IsAdmin(), routes.AgentRead)
		}

		room := api.Group("room")
		{
			room.GET("/divide", routes.RoomDivide)
			room.PUT("/resolve", middleware.IsAdmin(), routes.RoomResolve)
		}
	}

	route.Run()
}
