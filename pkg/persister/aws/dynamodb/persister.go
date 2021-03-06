package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/krzysztof-gzocha/pingor/pkg/config"
	"github.com/krzysztof-gzocha/pingor/pkg/persister/record"
	"github.com/pkg/errors"
)

// DynamoPutItem is interface used to simply inject and test AWS DynamoDB interaction
type DynamoPutItem interface {
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
}

// Persister is a struct responsible to store information about provided recordds
type Persister struct {
	client     DynamoPutItem
	config     config.DynamoDbPersister
	marshaller Marshaller
}

// NewPersister will return new DynamoDB persister
func NewPersister(client DynamoPutItem, config config.DynamoDbPersister) *Persister {
	return &Persister{
		client:     client,
		config:     config,
		marshaller: DynamoMarshaller{},
	}
}

// Persist will store provided result in AWS DynamoDB
func (p Persister) Persist(result record.Record) error {
	result.DeviceName = p.config.DeviceName

	marshalRes, err := p.marshaller.MarshalMap(result)
	if err != nil {
		return errors.Wrap(err, "unable to marshaller result into AWS object")
	}

	_, err = p.client.PutItem(&dynamodb.PutItemInput{
		Item:      marshalRes,
		TableName: aws.String(p.config.TableName),
	})

	return errors.Wrap(err, "unable to persist result in DynamoDB")
}
