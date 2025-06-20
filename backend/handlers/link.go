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

		// Fetch item_name from items table
		itemName := ""
		itemOut, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String("items"),
			Key: map[string]types.AttributeValue{
				"item_id": &types.AttributeValueMemberS{Value: req.ItemID},
			},
		})
		if err == nil {
			if n, ok := itemOut.Item["item_name"].(*types.AttributeValueMemberS); ok {
				itemName = n.Value
			}
		}

		// PutItem into links table
		_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String("links"),
			Item: map[string]types.AttributeValue{
				"fit_id":    &types.AttributeValueMemberS{Value: req.FitID},
				"item_id":   &types.AttributeValueMemberS{Value: req.ItemID},
				"item_name": &types.AttributeValueMemberS{Value: itemName},
			},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, echo.Map{"status": "linked"})
	}
}

// GetLinksForFitHandler handles GET /api/links/:fit_id
func GetLinksForFitHandler(client *dynamodb.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		fitID := c.Param("fit_id")
		if fitID == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing fit_id"})
		}
		// Query links table for all items for this fit
		out, err := client.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:              aws.String("links"),
			KeyConditionExpression: aws.String("fit_id = :fid"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":fid": &types.AttributeValueMemberS{Value: fitID},
			},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		links := []models.Link{}
		for _, v := range out.Items {
			fitID := ""
			if av, ok := v["fit_id"].(*types.AttributeValueMemberS); ok && av != nil {
				fitID = av.Value
			}
			itemID := ""
			if av, ok := v["item_id"].(*types.AttributeValueMemberS); ok && av != nil {
				itemID = av.Value
			}
			itemName := ""
			if av, ok := v["item_name"].(*types.AttributeValueMemberS); ok && av != nil {
				itemName = av.Value
			}
			link := models.Link{
				FitID:    fitID,
				ItemID:   itemID,
				ItemName: itemName,
			}
			links = append(links, link)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"links": links})
	}
}
