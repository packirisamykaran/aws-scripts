package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"net/http"
)

var tableName = "Collection-Watchlist" // Replace with your DynamoDB table name
var dynamoDBClient *dynamodb.DynamoDB

// Watchlist struct representing the DynamoDB item
type Watchlist struct {
	WalletAddress string   `json:"walletAddress"`
	Collection    []string `json:"collection"`
}

// Request struct representing the request body
type Request struct {
	WalletAddress  string `json:"walletAddress"`
	CollectionItem string `json:"collectionItem"`
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
	CollectionItem := requestBody.CollectionItem

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
		// Insert walletAddress and collectionItem into the Collection list and add to database
		err := insertNewWalletAddress(WalletAddress, CollectionItem)
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

		// Check if the collectionItem exists in the collection
		exists := false
		for i, item := range watchlist.Collection {
			if item == CollectionItem {
				exists = true
				// Remove the collectionItem from collection
				watchlist.Collection = append(watchlist.Collection[:i], watchlist.Collection[i+1:]...)
				break
			}
		}

		// Add the collectionItem to collection if it doesn't exist
		if !exists {
			watchlist.Collection = append(watchlist.Collection, CollectionItem)
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
		Body:       "Success",
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

// Function to insert a new walletAddress and collectionItem into DynamoDB
func insertNewWalletAddress(walletAddress, collectionItem string) error {
	watchlist := &Watchlist{
		WalletAddress: walletAddress,
		Collection:    []string{collectionItem},
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
