// +build unit

package dynamodb

import (
	"testing"

	"math"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/config"
	pkgMock "github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/krzysztof-gzocha/pingor/pkg/persister/record"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPersister(t *testing.T) {
	p := NewPersister(pkgMock.DynamoPutItemMock{}, config.DynamoDbPersister{})
	assert.NotNil(t, p)
	assert.IsType(t, &Persister{}, p)
}

func TestPersister_Persist(t *testing.T) {
	res := record.Record{}
	d := pkgMock.DynamoPutItemMock{}
	d.
		On("PutItem", mock.Anything).
		Once().
		Return(&dynamodb.PutItemOutput{}, nil)

	p := NewPersister(d, config.DynamoDbPersister{})
	err := p.Persist(res)

	assert.Nil(t, err)
	d.AssertExpectations(t)
}

func TestPersister_Persist_DynamoError(t *testing.T) {
	res := record.Record{}
	d := pkgMock.DynamoPutItemMock{}
	d.
		On("PutItem", mock.Anything).
		Once().
		Return(&dynamodb.PutItemOutput{}, errors.New("err"))

	p := NewPersister(d, config.DynamoDbPersister{})
	err := p.Persist(res)

	assert.Error(t, err)
	d.AssertExpectations(t)
}

func TestPersister_Persist_MarshalError(t *testing.T) {
	res := record.Record{
		CurrentResult: result.TimeResult{
			Result: result.Result{SuccessRate: float32(math.NaN())}},
	}
	d := pkgMock.DynamoPutItemMock{}

	marshaller := pkgMock.Marshaller{}
	marshaller.
		On("MarshalMap", mock.Anything).
		Once().
		Return(nil, errors.New("err"))

	p := NewPersister(d, config.DynamoDbPersister{})
	p.marshaller = marshaller
	err := p.Persist(res)

	assert.Error(t, err)
	marshaller.AssertExpectations(t)
}
