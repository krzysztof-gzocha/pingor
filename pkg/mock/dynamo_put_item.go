package mock

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

// DynamoPutItemMock is just a mock
type DynamoPutItemMock struct {
	mock.Mock
}

// PutItem is just a mock
func (m DynamoPutItemMock) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(input)

	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}
