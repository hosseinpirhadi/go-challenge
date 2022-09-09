package device

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorInvalidDeviceData       = "Invalid user data"
	ErrorFailedToFetchRecord     = "Failed to fetch record"
	ErrorFailedToUnmarshalRecord = "Failed to unmarshal record"
	ErrorDeviceAlreadyExists     = "Device already exists"
	ErrorCouldNotMarshalItem     = "Could not marshal item"
	ErrorCouldNotDynamoPutItem   = "Could not put item in dynamo"
	ErrorItemWasNotInTheDataBase = "Item was not in the DataBase"
	ErrorIdIsEmpty = "id field should be enterd"
	ErrorDeviceModelIsEmpty = "device model field should be enterd"
	ErrorNameIsEmpty = "name field should be enterd"
	ErrorNoteIsEmpty = "note field should be enterd"
	ErrorSerialIsEmpty = "serial field should be enterd"
)

type Device struct {
	Id          string `json:"id"`
	DeviceModel string `json:"deviceModel"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Serial      string `json:"serial"`
}

func FetchDevice(id, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Device, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(Device)
	
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	if item.Id == "" {
		return nil, errors.New(ErrorItemWasNotInTheDataBase)
	}

	return item, nil
}

func FetchDevices(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]Device, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]Device)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func CreateDevice(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Device,
	error,
) {
	var d Device

	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return nil, errors.New(ErrorInvalidDeviceData)
	}
	if d.Id == "" {
		return nil, errors.New(ErrorIdIsEmpty)
	}
	if d.DeviceModel == "" {
		return nil, errors.New(ErrorDeviceModelIsEmpty)
	}
	if d.Name == "" {
		return nil, errors.New(ErrorNameIsEmpty)
	}
	if d.Note == "" {
		return nil, errors.New(ErrorNoteIsEmpty)
	}
	if d.Serial == "" {
		return nil, errors.New(ErrorSerialIsEmpty)
	}
	currentDevice, _ := FetchDevice(d.Id, tableName, dynaClient)
	if currentDevice != nil && len(currentDevice.Id) != 0 {
		return nil, errors.New(ErrorDeviceAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(d)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &d, nil
}
