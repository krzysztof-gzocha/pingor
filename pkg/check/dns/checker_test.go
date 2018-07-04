// +build unit

package dns

import (
	"testing"

	"time"

	"context"

	internalMock "github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDNSChecker(t *testing.T) {
	successDnsMock := new(dnsMock)
	result := Result{Time: time.Second}
	successDnsMock.
		On("ResolveHost", "wp.pl").
		Once().
		Return(result, nil)

	checker := NewChecker(&internalMock.Logger{}, successDnsMock, "wp.pl")
	assert.Equal(t, checker.dns, successDnsMock)
	assert.Len(t, checker.hosts, 1)
}

func TestDNSChecker_Check(t *testing.T) {
	successDnsMock := new(dnsMock)
	unsuccessfulDnsMock := new(dnsMock)
	dnsResult := Result{Time: time.Second}
	hosts := []string{"wp.pl", "onet.pl"}
	successDnsMock.
		On("ResolveHost", mock.Anything).
		Times(len(hosts)).
		Return(dnsResult, nil)
	unsuccessfulDnsMock.
		On("ResolveHost", mock.Anything).
		Times(len(hosts)).
		Return(dnsResult, errors.New("error"))

	logger := &internalMock.Logger{}
	logger.On("WithField", "host", "wp.pl")
	logger.On("Debugf", mock.Anything, mock.Anything)
	logger.On("WithField", "host", "onet.pl")
	logger.On("Debugf", mock.Anything, mock.Anything)
	logger.On("WithField", "successRate", mock.Anything)
	logger.On("Debugf", mock.Anything, mock.Anything)

	checker := NewChecker(logger, successDnsMock, hosts...)
	result := checker.Check(context.TODO())
	assert.True(t, result.IsSuccess())
	assert.Equal(t, result.GetTime(), time.Second)
	assert.Len(t, result.GetSubResults(), 2)
	assert.True(t, successDnsMock.AssertExpectations(t))
	logger.AssertExpectations(t)

	logger = &internalMock.Logger{}
	logger.On("WithField", "host", "wp.pl")
	logger.On("Debugf", mock.Anything, mock.Anything)
	logger.On("WithField", "host", "onet.pl")
	logger.On("Debugf", mock.Anything, mock.Anything)
	logger.On("WithField", "successRate", mock.Anything)
	logger.On("Debugf", mock.Anything, mock.Anything)

	checker = NewChecker(logger, unsuccessfulDnsMock, hosts...)
	result = checker.Check(context.TODO())
	assert.False(t, result.IsSuccess())
	assert.Equal(t, result.GetTime(), time.Second)
	assert.Len(t, result.GetSubResults(), 2)
	assert.Equal(t, "Checking DNS with 2 hosts", result.GetMessage())
	assert.True(t, unsuccessfulDnsMock.AssertExpectations(t))
	logger.AssertExpectations(t)

	logger = &internalMock.Logger{}
	checker = NewChecker(logger, successDnsMock)
	result = checker.Check(context.TODO())
	assert.False(t, result.IsSuccess())
	assert.Zero(t, result.GetSubResults())
	assert.Zero(t, result.GetSuccessRate())
	assert.Zero(t, result.GetTime())
	logger.AssertExpectations(t)
}

type dnsMock struct {
	mock.Mock
	result Result
}

func (m dnsMock) ResolveHost(host string) (Result, error) {
	args := m.Called(host)

	return args.Get(0).(Result), args.Error(1)
}
