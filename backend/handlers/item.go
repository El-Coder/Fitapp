package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"fitapp-backend/models"
)

// CreateItemHandler handles POST /api/items
func CreateItemHandler(client *dynamodb.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var item models.Item
		if err := c.Bind(&item); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		if item.ItemID == "" {
			item.ItemID = uuid.New().String()
		}
		_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String("items"),
			Item: map[string]types.AttributeValue{
				"item_id":   &types.AttributeValueMemberS{Value: item.ItemID},
				"item_name": &types.AttributeValueMemberS{Value: item.ItemName},
			},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusCreated, item)
	}
}

// GetItemsHandler handles GET /api/items
func GetItemsHandler(client *dynamodb.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		out, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String("items"),
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		items := []models.Item{}
		for _, v := range out.Items {
			item := models.Item{
				ItemID:   v["item_id"].(*types.AttributeValueMemberS).Value,
				ItemName: v["item_name"].(*types.AttributeValueMemberS).Value,
			}
			items = append(items, item)
		}
		return c.JSON(http.StatusOK, items)
	}
}
