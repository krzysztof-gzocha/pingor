package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Marshaller is an interface helping with testing dynamo db persister
type Marshaller interface {
	MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error)
}

// DynamoMarshaller will re-use dynamodbattribute.MarshalMap function. It's just hidden within the struct for easier testing
type DynamoMarshaller struct{}

// MarshalMap will call dynamodbattribute.MarshalMap
func (m DynamoMarshaller) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}
