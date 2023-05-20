package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var tableName = "Collection-Watchlist" // Replace with your DynamoDB table name

// Request struct to parse the incoming API Gateway request
type Request struct {
	WalletAddress string `json:"walletAddress"`
	Value         string `json:"value"`
}

// Response struct for the API Gateway response
type Response struct {
	Message string `json:"message"`
}

// Handler function to process the Lambda event
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the request body
	// var req Request
	// if err := request.UnmarshalJSON(&req); err != nil {
	// 	return Response{}, fmt.Errorf("failed to parse request body: %v", err)

	WalletAddress := "b"
	Value := "sui"
	// }

	// Create a new DynamoDB session
	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string("error"),
		}, nil
	}

	// Create a new DynamoDB client
	svc := dynamodb.New(sess)

	// Get the existing collection for the given wallet address
	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"walletAddress": {
				S: aws.String(WalletAddress),
			},
		},
	}

	getResult, err := svc.GetItem(getInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string("error"),
		}, nil
	}

	// Update the collection based on the request value
	updateExpression := ""
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{}
	if getResult.Item == nil {
		// If the collection doesn't exist, create a new collection with the request value
		updateExpression = "SET #col = :val"
		expressionAttributeValues[":val"] = &dynamodb.AttributeValue{
			SS: []*string{aws.String(Value)},
		}
	} else {
		// If the collection exists, add the request value if it doesn't exist, or remove it if it already exists
		collection := getResult.Item["collection"].SS
		isValueExists := false
		for _, v := range collection {
			if *v == Value {
				isValueExists = true
				break
			}
		}

		if isValueExists {
			updateExpression = "DELETE #col :val"
			expressionAttributeValues[":val"] = &dynamodb.AttributeValue{
				SS: []*string{aws.String(Value)},
			}
		} else {
			updateExpression = "ADD #col :val"
			expressionAttributeValues[":val"] = &dynamodb.AttributeValue{
				SS: []*string{aws.String(Value)},
			}
		}
	}

	// Update the collection for the given wallet address
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"walletAddress": {
				S: aws.String(WalletAddress),
			},
		},
		UpdateExpression: aws.String(updateExpression),
		ExpressionAttributeNames: map[string]*string{
			"#col": aws.String("collection"),
		},
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              aws.String("ALL_NEW"),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string("error"),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string("updated"),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
