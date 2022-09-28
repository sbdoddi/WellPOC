package employee

import(
	"github.com/sbdoddi/golangCrud/pkg/validation"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var(
	Deleted = " Deleted Successfully"
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord = "failed to fetch record"
	ErrorInvalidUserData = "invalid user data"
	ErrorInvalidEmail = "invalid email"
	ErrorCouldNotMarshalItem = "could not marshal item"
	ErrorCouldNotDeleteItem = "could not delete item"
	ErrorCouldNotDynamoPutItem = "could not dynamo put item"
	ErrorUserAlreadyExists = "user.User already exists"
	ErrorUserDoesNotExist = "user.User does not exist"
)

type Employee struct{
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}

func FetchEmployee(email string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Employee, error) {
	input := &dynamodb.GetItemInput{
		Key:map[string]*dynamodb.AttributeValue{
			"email":{
				S:aws.String(email),
			},
		},
		TableName:aws.String(tableName),
	}

	
	result, err := dynaClient.GetItem(input)
	if err!= nil{
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(Employee)

	err1 := dynamodbattribute.UnmarshalMap(result.Item, item)

	if err1!= nil{
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchEmployees(tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*[]Employee, error){
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err!= nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]Employee)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func CreateEmployee(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*Employee, error)  {

	var e Employee

	err := json.Unmarshal([]byte(req.Body), &e)
	if err != nil{
		return nil, errors.New(ErrorInvalidUserData)
	}
	if validation.IsEmailValid(e.Email) == false{
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentEmployee, _ := FetchEmployee(e.Email, tableName, dynaClient)
	if currentEmployee != nil && len(currentEmployee.Email) != 0{
		return nil, errors.New(ErrorUserAlreadyExists)
	}
	av ,err := dynamodbattribute.MarshalMap(e)
	if err != nil{
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}
	input := &dynamodb.PutItemInput{
		Item:av,
		TableName:aws.String(tableName),
	}

	_, err1 := dynaClient.PutItem(input)
	if err1 != nil{
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &e, nil
	
}

func UpdateEmployee(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Employee, error) {
	var e Employee
	err := json.Unmarshal([]byte(req.Body), &e)
	if err!= nil{
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	if validation.IsEmailValid(e.Email) == false{
		return nil, errors.New(ErrorUserAlreadyExists)
	}
	/*currentEmployee, err := FetchEmployee(e.Email, tableName, dynaClient)
	if currentEmployee != nil && len(currentEmployee.Email) != 0{
		return nil, errors.New(ErrorUserAlreadyExists)
	} */
	av,err := dynamodbattribute.MarshalMap(e)
	if err!= nil{
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}
	input := &dynamodb.PutItemInput{
		Item:av,
		TableName:aws.String(tableName),
	}

	_, err1 := dynaClient.PutItem(input)
	if err1 != nil{
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &e, nil
}

func DeleteEmployee(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Employee, error) {
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key:map[string]*dynamodb.AttributeValue{
			"email":{
				S:aws.String(email),
			},
		},
		TableName:aws.String(tableName),
	}
	_, err := dynaClient.DeleteItem(input)
	if err!= nil{
		return nil, errors.New(ErrorCouldNotDeleteItem)
	}
	return nil, errors.New(Deleted)
}