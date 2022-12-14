// This package is the core of the project
package main

import (
	"os"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/hosseinpirhadi/challenge/pkg/handlers"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)


func main() {
	region := os.Getenv("AWS_REGION")
	// create a session for connectiong to dunamodb
	awsSession, err := session.NewSession(&aws.Config{	
		Region: aws.String(region)})

	if err != nil {
		return
	}

	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

//dynamodb table name
const tableName = "device-table-challenge"

//handler gets all requests from AWS Api gateway
func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	switch req.HTTPMethod {
	case "GET":
		return handlers.GetDevice(req, tableName, dynaClient)
	case "POST":
		return handlers.CreateDevice(req, tableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}
}
