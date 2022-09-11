package handlers

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

//apiResponse will get responses from handlers and device methods and create a json to send to client
func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}
