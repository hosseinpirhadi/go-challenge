package handlers

import (
	"net/http"
	"strings"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/hosseinpirhadi/challenge/pkg/device"
)

var ErrorMethodNotAllowed = "method not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

//GetDevice handle Get request
//if the client request a specific device id this function will call FetchDevice
//if the clieant request whitout a specific id this function will call FetchDevices
func GetDevice(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	path := req.Path
	splitedPath := strings.Split(path, "/")
	id := splitedPath[len(splitedPath) - 1]

	if len(id) > 0 {
		result, err := device.FetchDevice(id, tableName, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{
				aws.String(err.Error()),
			})
		}
		return apiResponse(http.StatusOK, result)
	}

	result, err := device.FetchDevices(tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)
}

//CreateDecive handle Post request
func CreateDevice(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := device.CreateDevice(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusCreated, result)
}

//UnhandledMethod will handle requests other than Get and Post
func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
