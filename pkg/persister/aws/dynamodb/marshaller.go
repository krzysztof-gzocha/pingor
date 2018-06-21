package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type MarshallerInterface interface {
	MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error)
}

type Marshaller struct{}

func (m Marshaller) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}
