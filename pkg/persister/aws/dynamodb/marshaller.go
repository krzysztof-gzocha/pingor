package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// MarshallerInterface is an interface helping with testing dynamo db persister
type MarshallerInterface interface {
	MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error)
}

// Marshaller will re-use dynamodbattribute.MarshalMap function. It's just hidden within the struct for easier testing
type Marshaller struct{}

// MarshalMap will call dynamodbattribute.MarshalMap
func (m Marshaller) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}
