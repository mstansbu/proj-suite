package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mstansbu/tic-tac-toe/templates"
)

func main() {
	server := gin.Default()

	server.StaticFile("main.css", "./templates/main.css")

	server.GET("/", serveHome)

	server.Run(":3000")
}

func serveHome(c *gin.Context) {
	templates.Layout().Render(c.Request.Context(), c.Writer)
}
