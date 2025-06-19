package main

import (
	"context"
	"fitapp-backend/db"
	"fitapp-backend/routes"
	"log"

	"github.com/labstack/echo/v4/middleware"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
)

func EnsureFitsTable(client *dynamodb.Client) {
	tb := db.TableBasics{
		DynamoDbClient: client,
		TableName:      "fits",
	}
	_, err := tb.CreateFitsTable(context.TODO())
	if err != nil {
		log.Fatalf("Failed to ensure 'fits' table: %v", err)
	}
}

func EnsureItemsTable(client *dynamodb.Client) {
	tb := db.TableBasics{
		DynamoDbClient: client,
		TableName:      "items",
	}
	_, err := tb.CreateItemsTable(context.TODO())
	if err != nil {
		log.Fatalf("Failed to ensure 'items' table: %v", err)
	}
}

func main() {
	client := db.CreateLocalClient()
	db.WaitForDynamoReady("http://dynamodb-local:8000")
	EnsureFitsTable(client)
	EnsureItemsTable(client)

	e := echo.New()
	e.Use(middleware.CORS())

	routes.RegisterRoutes(e, client)
	e.Logger.Fatal(e.Start(":8080"))
}
