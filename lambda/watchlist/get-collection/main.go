package main

import (
	"context"
	"encoding/json"

	"log"
	"net/http"

	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var tableName = "collection-watchlist-solana" // Replace with your DynamoDB table name
var dynamoDBClient *dynamodb.DynamoDB

// Watchlist struct representing the DynamoDB item
type Watchlist struct {
	WalletAddress    string   `json:"walletAddress"`
	CollectionIDList []string `json:"collectionIDList"`
}

// Handler function to process the Lambda event
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the request body
	// walletAddress := "b"
	walletAddress := request.QueryStringParameters["walletAddress"]
	// Get the collection by walletAddress
	result, err := getCollectionIDList(walletAddress)
	if err != nil {
		log.Printf("Error getting collectionIDList: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	// If walletAddress doesn't exist, create a new row with an empty collection
	if result == nil {
		if err := createEmptyWatchlist(walletAddress); err != nil {
			log.Printf("Error creating empty watchlist: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
			}, nil
		}
	}

	// Build the response
	responseBody, err := json.Marshal(result.CollectionIDList)
	if err != nil {
		log.Printf("Error marshaling response body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	var responseList []string
	if err := json.Unmarshal(responseBody, &responseList); err != nil {
		log.Printf("Error unmarshaling response body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	responseString := strings.Join(responseList, ",")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       responseString,
	}, nil
}

// Function to get the collection by walletAddress from DynamoDB
func getCollectionIDList(walletAddress string) (*Watchlist, error) {
	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"walletAddress": {
				S: aws.String(walletAddress),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return &Watchlist{
			WalletAddress:    walletAddress,
			CollectionIDList: []string{},
		}, nil
	}

	watchlist := &Watchlist{}
	err = dynamodbattribute.UnmarshalMap(result.Item, watchlist)
	if err != nil {
		return nil, err
	}

	return watchlist, nil
}

// Function to create an empty watchlist with the given walletAddress
func createEmptyWatchlist(walletAddress string) error {
	watchlist := &Watchlist{
		WalletAddress:    walletAddress,
		CollectionIDList: []string{},
	}

	item, err := dynamodbattribute.MarshalMap(watchlist)
	if err != nil {
		return err
	}

	_, err = dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Create a new DynamoDB session
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("Failed to create DynamoDB session: %v", err)
	}

	// Create a new DynamoDB client
	dynamoDBClient = dynamodb.New(sess)

	lambda.Start(Handler)
}
