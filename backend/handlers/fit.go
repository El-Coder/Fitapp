package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"fitapp-backend/models"
	"github.com/labstack/echo/v4"
)

// CreateFitHandler handles POST /api/fits
func CreateFitHandler(client *dynamodb.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var fit models.Fit
		if err := c.Bind(&fit); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		if fit.FitID == "" {
			fit.FitID = uuid.New().String()
		}
		_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String("fits"),
			Item: map[string]types.AttributeValue{
				"fit_id":   &types.AttributeValueMemberS{Value: fit.FitID},
				"fit_name": &types.AttributeValueMemberS{Value: fit.FitName},
			},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusCreated, fit)
	}
}

// GetFitsHandler handles GET /api/fits
func GetFitsHandler(client *dynamodb.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		out, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String("fits"),
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		fits := []models.Fit{}
		for _, v := range out.Items {
			fit := models.Fit{
				FitID:   v["fit_id"].(*types.AttributeValueMemberS).Value,
				FitName: v["fit_name"].(*types.AttributeValueMemberS).Value,
			}
			fits = append(fits, fit)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"fits": fits})
	}
}
