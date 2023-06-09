package main

import (
	"Backend/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)
func main() {
	// Echo instance
	e := echo.New()
	
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	
	// controller.Processor.Query = "Apa Mata Kuliah paling seru di semester 4"

	// Routes
	e.POST("/stimaGPT", controller.CreateQnA)
	e.GET("/stimaGPT", controller.GetListQnA)
	e.DELETE("/stimaGPT/delete", controller.DelQnA)

	e.POST("/stimaGPT/User", controller.CreateUser)
	e.GET("/stimaGPT/User", controller.GetUser)
	
	e.POST("/stimaGPT/history", controller.CreateHistory)
	e.GET("/stimaGPT/history", controller.GetHistory)
	e.DELETE("/stimaGPT/history", controller.DeleteHistory)

	e.POST("/stimaGPT/chat", controller.GetAnswers)
	e.GET("/stimaGPT/chat", controller.GetChatHistory)
	// Start server
	e.Logger.Fatal(e.Start(":5000"))
}




