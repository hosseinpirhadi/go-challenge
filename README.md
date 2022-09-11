![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/hosseinpirhadi/go-challenge)


# Go Challenge: Simple Restful API on AWS

This is a simple Restful API on AWS using the following tech stack:

- [Serverless Framework](https://serverless.com/)
- [Go language](https://golang.org/)
- AWS API Gateway
- AWS Lambda
- AWS DynamoDB

The API accepts the following JSON requests and produces the corresponding HTTP responses:

## Request 1 (HTTP POST)

URL: `https://<api-gateway-url>/api/devices`
```json
{
  "id": "/devices/id1",
  "deviceModel": "/devicemodels/id1",
  "name": "Sensor",
  "note": "Testing a sensor.",
  "serial": "A020000102"
}
```
#### Response 1 - Success
`HTTP 201 Created`

#### Response 1 - Failure 1
`HTTP 400 Bad Request` (If any of the payload fields are missing)

The response body has a descriptive error message for the client to be able to detect the problem.

#### Response 1 - Failure 2
`HTTP 500 Internal Server Error` (If any exceptional situation occurs on the server side)

## Request 2 (HTTP GET)

URL: `https://<api-gateway-url>/api/devices/{id}`

Example: `https://api123.amazonaws.com/api/devices/id1`

#### Response 2 - Success
`HTTP 200 OK`
```json
{
  "id": "/devices/id1",
  "deviceModel": "/devicemodels/id1",
  "name": "Sensor",
  "note": "Testing a sensor.",
  "serial": "A020000102"
}
```

#### Response 2 - Failure 1
`HTTP 404 Not Found` (If the request-id does not exist)

#### Response 2 - Failure 2
`HTTP 500 Internal Server Error` (If any exceptional situation occurs on the server side)

## Usage

### Prerequisites
- [The Go Programming Language](https://go.dev/doc/install)
- [Get Docker](https://docs.docker.com/get-docker/)
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
- use `go mod tidy` in the root directory of project to download dependencies

### Testing
- Build the project using `sam build` command
- To run the project locally, use `sam local start-api` command
- It can be tested using PostMan

![testing challenge via postman](./testing.gif?raw=true)

<p><img align="right" alt="git" src"https://github.com/hosseinpirhadi/go-challenge/blob/main/testing.gif" width="500" height="320" /> </p>


Use [this link](https://8xsnc3ax5l.execute-api.us-east-1.amazonaws.com/staging/) to test project on AWS like above example.

## pkg

### Package device

#### variables
Possible errors are listed below.
```go
var (
    ErrorInvalidDeviceData       = "Invalid user data"
    ErrorFailedToFetchRecord     = "Failed to fetch record"
    ErrorFailedToUnmarshalRecord = "Failed to unmarshal record"
    ErrorDeviceAlreadyExists     = "Device already exists"
    ErrorCouldNotMarshalItem     = "Could not marshal item"
    ErrorCouldNotDynamoPutItem   = "Could not put item in dynamo"
    ErrorItemWasNotInTheDataBase = "Item was not in the DataBase"
)
```

#### _func_ FetchDevices
Obtain all devices in the database.
```go
func FetchDevices(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]Device, error)
```

#### _type_ Device
```go
type Device struct {
    Id          string `json:"id"`
    DeviceModel string `json:"deviceModel"`
    Name        string `json:"name"`
    Note        string `json:"note"`
    Serial      string `json:"serial"`
}
```

#### _func_ CreateDevice
Generates a new device with the aforementioned features.
```go
func CreateDevice(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
    *Device,
    error,
)
```

#### _func_ FetchDevice
Get a particular device from the database.
```go
func FetchDevice(id, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Device, error)
```

### Package handlers

#### Variables
When a user requests a method that does not exist, the error `"method not allowed"` appears.
```go
var ErrorMethodNotAllowed = "method not allowed"
```

#### _func_ CreateDevice 
Handle the situation once the user requests the creation of a new device.
```go
func CreateDevice(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error)
```

#### _func_ GetDevice 
Whenever the user requests a specific device, this method invokes the "FetchDevice" operation.
```go
func GetDevice(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error)
```

#### _func_ UnhandledMethod
If a user requests a method that does not exist, an error will be sent to the API response.
```go
func UnhandledMethod() (*events.APIGatewayProxyResponse, error)
```

#### _type_ ErrorBody
```go
type ErrorBody struct {
    ErrorMsg *string `json:"error,omitempty"`
}
```
