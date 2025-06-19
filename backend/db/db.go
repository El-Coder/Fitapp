package db

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func CreateLocalClient() *dynamodb.Client {
	endpoint := "http://dynamodb-local:8000"
	region := "us-east-1"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "DUMMYIDEXAMPLE",
				SecretAccessKey: "DUMMYEXAMPLEKEY",
				SessionToken:    "dummy",
				Source:          "Local dummy credentials",
			},
		}),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			}),
		),
	)

	if err != nil {
		panic("failed to load SDK config: " + err.Error())
	}

	return dynamodb.NewFromConfig(cfg)
}

func WaitForDynamoReady(url string) {
	for i := 0; i < 10; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 400 {
			log.Println("DynamoDB is ready.")
			return
		}
		log.Println("Waiting for DynamoDB to be ready...")
		time.Sleep(2 * time.Second)
	}
	log.Fatal("DynamoDB not reachable after multiple attempts")
}

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

// CreateItemsTable creates the "items" table if it doesn't already exist.
func (basics TableBasics) CreateItemsTable(ctx context.Context) (*types.TableDescription, error) {
	// Check if the table already exists
	out, err := basics.DynamoDbClient.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		log.Printf("Failed to list tables: %v", err)
		return nil, err
	}
	for _, table := range out.TableNames {
		if table == basics.TableName {
			log.Printf("Table %v already exists. Skipping creation.", basics.TableName)
			return nil, nil
		}
	}

	table, err := basics.DynamoDbClient.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("item_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("item_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(basics.TableName),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Reason: %v", basics.TableName, err)
		return nil, err
	}

	waiter := dynamodb.NewTableExistsWaiter(basics.DynamoDbClient)
	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(basics.TableName),
	}, 5*time.Minute)
	if err != nil {
		log.Printf("Table waiter failed for %v: %v", basics.TableName, err)
		return nil, err
	}

	log.Printf("Created table: %v", basics.TableName)
	return table.TableDescription, nil
}

// / CreateFitsTable creates the "fits" table if it doesn't already exist.
func (basics TableBasics) CreateFitsTable(ctx context.Context) (*types.TableDescription, error) {
	// Check if the table already exists
	out, err := basics.DynamoDbClient.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		log.Printf("Failed to list tables: %v", err)
		return nil, err
	}
	for _, table := range out.TableNames {
		if table == basics.TableName {
			log.Printf("Table %v already exists. Skipping creation.", basics.TableName)
			return nil, nil
		}
	}

	// Create the table if it doesn't exist
	table, err := basics.DynamoDbClient.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("fit_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("fit_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(basics.TableName),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Reason: %v", basics.TableName, err)
		return nil, err
	}

	waiter := dynamodb.NewTableExistsWaiter(basics.DynamoDbClient)
	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(basics.TableName),
	}, 5*time.Minute)
	if err != nil {
		log.Printf("Table waiter failed for %v: %v", basics.TableName, err)
		return nil, err
	}

	log.Printf("Created table: %v", basics.TableName)
	return table.TableDescription, nil
}
