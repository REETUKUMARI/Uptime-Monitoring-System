package main

import (
	"website_sc/controllers"
	"website_sc/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	models.ConnectDataBase()
	go controllers.Checklink()

	r.GET("/urls/", controllers.GetUrls)
	r.GET("/urls/:id", controllers.GetUrl)
	r.POST("/urls", controllers.CreateUrl)
	r.PATCH("urls/:id", controllers.Updateurl)
	r.DELETE("urls/:id", controllers.Deleteurl)

	r.Run(":8080")
}
