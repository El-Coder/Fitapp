package routes

import (
	"fitapp-backend/handlers"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, client *dynamodb.Client) {
	e.GET("/ping-db", handlers.HealthCheckHandler(client))
	e.POST("/api/link", handlers.LinkItemHandler(client))

	// New item endpoints
	e.POST("/api/items", handlers.CreateItemHandler(client))
	e.GET("/api/items", handlers.GetItemsHandler(client))

	// New fit endpoints
	e.POST("/api/fits", handlers.CreateFitHandler(client))
	e.GET("/api/fits", handlers.GetFitsHandler(client))

	// New links endpoint
	e.GET("/api/links/:fit_id", handlers.GetLinksForFitHandler(client))
}
