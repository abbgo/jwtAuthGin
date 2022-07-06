package routes

import (
	"jwtAuthGin/controllers"
	"jwtAuthGin/middlewares"

	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {

	routes := gin.Default()

	api := routes.Group("/api")
	{
		api.POST("/login", controllers.Login)                // ulanyjy login bolanda acces token yasalyar
		api.POST("/user/register", controllers.RegisterUser) // ulanyjy registrasiya bolyar
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", controllers.Ping)
		}
		api.GET("/refresh", controllers.Refresh)
	}

	return routes

}
