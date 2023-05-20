package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(handler)
}

type MyEvent struct {
	WalletAddress string `json:"walletAddress"`
}

type CollectionResponse struct {
	Collection []string `json:"collection"`
}

// arn:aws:dynamodb:ap-southeast-2:120657039516:table/Collection-Watchlist

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// walletAddress := request.QueryStringParameters["walletAddress"]

	// err := json.Unmarshal([]byte(request.Body), &bodyData)

	svc := dynamodb.New(session.New())

	walletAddress := "a"

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Collection-Watchlist"), // Replace with your table name
		Key: map[string]*dynamodb.AttributeValue{
			"walletAddress": {
				S: aws.String(walletAddress),
			},
		},
	}

	// Retrieve the item from DynamoDB
	result, err := svc.GetItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Process the result
	if result.Item == nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("Item not found")
	}

	// Extract the collection attribute value
	collectionAttr := result.Item["collection"]
	if collectionAttr == nil || collectionAttr.SS == nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("Collection attribute not found")
	}

	// Convert []*string to []string
	collection := make([]string, 0)
	for _, s := range collectionAttr.SS {
		collection = append(collection, *s)
	}

	// Create the response
	response := CollectionResponse{
		Collection: collection,
	}
	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Return the response
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil

}

type GetQuery struct {
	walletAddress string
}
