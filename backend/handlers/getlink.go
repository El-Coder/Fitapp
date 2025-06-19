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
			linked := []string{}
			if v, ok := item["fit_id"].(*types.AttributeValueMemberS); ok {
				fitID = v.Value
			}
			if v, ok := item["linked_item_ids"].(*types.AttributeValueMemberSS); ok {
				linked = v.Value
			}
			results = append(results, map[string]interface{}{
				"fit_id":          fitID,
				"linked_item_ids": linked,
			})
		}

		return c.JSON(http.StatusOK, echo.Map{"fits": results})
	}
}
