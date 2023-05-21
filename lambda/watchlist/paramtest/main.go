package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Read URL parameter from event.QueryStringParameters
	walletAddress := request.QueryStringParameters["walletAddress"]

	// Do something with the walletAddress
	responseBody := map[string]string{
		"walletAddress": walletAddress,
	}

	// Convert response body to JSON
	responseJSON, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error marshaling response body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJSON),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
