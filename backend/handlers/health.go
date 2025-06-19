package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
)

func HealthCheckHandler(client *dynamodb.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		out, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, out.TableNames)
	}
}
