package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

// Request struct representing the request body
type Request struct {
	WalletAddress string `json:"walletAddress"`
	CollectionID  string `json:"collectionID"`
}

// Handler function to process the Lambda event
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the request body

	// Read request body
	var requestBody Request
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Printf("Error unmarshaling request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request body",
		}, nil
	}

	// Access the fields of the request body
	WalletAddress := requestBody.WalletAddress
	CollectionID := requestBody.CollectionID

	// Check if the walletAddress exists
	exists, err := checkWalletAddressExists(WalletAddress)
	if err != nil {
		log.Printf("Error checking wallet address: %v", err)
		return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
			},
			nil
	}

	if !exists {
		// Insert walletAddress and collectionID into the Collection list and add to database
		err := insertNewWalletAddress(WalletAddress, CollectionID)
		if err != nil {
			log.Printf("Error inserting wallet address: %v", err)
			return events.APIGatewayProxyResponse{
					StatusCode: http.StatusInternalServerError,
					Body:       http.StatusText(http.StatusInternalServerError),
				},
				nil
		}
	} else {
		// Get the current collection for the given walletAddress
		watchlist, err := getWatchlist(WalletAddress)
		if err != nil {
			log.Printf("Error getting watchlist: %v", err)
			return events.APIGatewayProxyResponse{
					StatusCode: http.StatusInternalServerError,
					Body:       http.StatusText(http.StatusInternalServerError),
				},
				nil
		}

		// Check if the collectionID exists in the collection
		exists := false
		for i, item := range watchlist.CollectionIDList {
			if item == CollectionID {
				exists = true
				// Remove the collectionID from collection
				watchlist.CollectionIDList = append(watchlist.CollectionIDList[:i], watchlist.CollectionIDList[i+1:]...)
				break
			}
		}

		// Add the collectionID to collection if it doesn't exist
		if !exists {
			watchlist.CollectionIDList = append(watchlist.CollectionIDList, CollectionID)
		}

		// Save the updated watchlist to DynamoDB
		err = saveWatchlist(watchlist)
		if err != nil {
			log.Printf("Error saving watchlist: %v", err)
			return events.APIGatewayProxyResponse{
					StatusCode: http.StatusInternalServerError,
					Body:       http.StatusText(http.StatusInternalServerError),
				},
				nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "updated",
	}, nil
}

// Function to check if the walletAddress exists in DynamoDB
func checkWalletAddressExists(walletAddress string) (bool, error) {
	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"walletAddress": {
				S: aws.String(walletAddress),
			},
		},
	})
	if err != nil {
		return false, err
	}

	return result.Item != nil, nil
}

// Function to insert a new walletAddress and collectionID into DynamoDB
func insertNewWalletAddress(walletAddress, collectionID string) error {
	watchlist := &Watchlist{
		WalletAddress:    walletAddress,
		CollectionIDList: []string{collectionID},
	}

	item, err := dynamodbattribute.MarshalMap(watchlist)
	if err != nil {
		return err
	}

	_, err = dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

// Function to get the watchlist by walletAddress from DynamoDB
func getWatchlist(walletAddress string) (*Watchlist, error) {
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
		return nil, fmt.Errorf("Watchlist not found for walletAddress: %s", walletAddress)
	}

	watchlist := &Watchlist{}
	err = dynamodbattribute.UnmarshalMap(result.Item, watchlist)
	if err != nil {
		return nil, err
	}

	return watchlist, nil
}

// Function to save the watchlist to DynamoDB
func saveWatchlist(watchlist *Watchlist) error {
	item, err := dynamodbattribute.MarshalMap(watchlist)
	if err != nil {
		return err
	}

	_, err = dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	dynamoDBClient = dynamodb.New(sess)

	lambda.Start(Handler)
}
