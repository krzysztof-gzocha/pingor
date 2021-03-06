package mock

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

// DynamoMarshaller is just a mock
type Marshaller struct {
	mock.Mock
}

// MarshalMap is just a mock
func (m Marshaller) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	args := m.Called(in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]*dynamodb.AttributeValue), args.Error(1)
}
