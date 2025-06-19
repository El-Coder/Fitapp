package handlers

import (
	"context"
	"fitapp-backend/models"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/labstack/echo/v4"
)

func LinkItemHandler(client *dynamodb.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.LinkRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
		}

		// DynamoDB logic
		_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName: aws.String("fits"),
			Key: map[string]types.AttributeValue{
				"fit_id": &types.AttributeValueMemberS{Value: req.FitID},
			},
			UpdateExpression: aws.String("ADD linked_item_ids :i"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":i": &types.AttributeValueMemberSS{Value: []string{req.ItemID}},
			},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, echo.Map{"status": "linked"})
	}
}
