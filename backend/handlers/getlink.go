package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/labstack/echo/v4"
)

func GetAllLinkedItemsHandler(client *dynamodb.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		out, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName:            aws.String("fits"),
			ProjectionExpression: aws.String("fit_id, linked_item_ids"),
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		results := []map[string]interface{}{}
		for _, item := range out.Items {
			fitID := ""
			fitName := ""
			linked := []string{}
			linkedItems := []map[string]string{}

			if v, ok := item["fit_id"].(*types.AttributeValueMemberS); ok {
				fitID = v.Value
				// Fetch fit name from fits table
				fitOut, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
					TableName: aws.String("fits"),
					Key: map[string]types.AttributeValue{
						"fit_id": &types.AttributeValueMemberS{Value: fitID},
					},
				})
				if err == nil {
					if n, ok := fitOut.Item["fit_name"].(*types.AttributeValueMemberS); ok {
						fitName = n.Value
					}
				}
			}
			if v, ok := item["linked_item_ids"].(*types.AttributeValueMemberSS); ok {
				linked = v.Value
				for _, itemID := range linked {
					itemName := ""
					itemOut, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
						TableName: aws.String("items"),
						Key: map[string]types.AttributeValue{
							"item_id": &types.AttributeValueMemberS{Value: itemID},
						},
					})
					if err == nil {
						if n, ok := itemOut.Item["item_name"].(*types.AttributeValueMemberS); ok {
							itemName = n.Value
						}
					}
					linkedItems = append(linkedItems, map[string]string{
						"item_id":   itemID,
						"item_name": itemName,
					})
				}
			}
			results = append(results, map[string]interface{}{
				"fit_id":          fitID,
				"fit_name":        fitName,
				"linked_items":    linkedItems,
			})
		}

		return c.JSON(http.StatusOK, echo.Map{"fits": results})
	}
}
