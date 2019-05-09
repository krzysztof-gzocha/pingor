// +build unit

package dns

import (
	"context"
	"testing"
	"time"

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
}

func TestDNSChecker_Check_Success(t *testing.T) {
	successDnsMock := new(dnsMock)
	// unsuccessfulDnsMock := new(dnsMock)
	dnsResult := Result{Time: time.Second}
	host := "wp.pl"
	successDnsMock.
		On("ResolveHost", mock.Anything).
		Once().
		Return(dnsResult, nil)

	logger := &internalMock.Logger{}
	logger.On("WithField", "host", "wp.pl")
	logger.On("Debugf", mock.Anything, mock.Anything)

	checker := NewChecker(logger, successDnsMock, host)
	result := checker.Check(context.TODO())
	assert.True(t, result.IsSuccess())
	assert.Equal(t, result.GetTime(), time.Second)
	assert.True(t, successDnsMock.AssertExpectations(t))
	logger.AssertExpectations(t)
}

func TestDNSChecker_Check_Failure(t *testing.T) {
	dnsMock := new(dnsMock)
	dnsResult := Result{Time: time.Second}
	host := "wp.pl"
	dnsMock.
		On("ResolveHost", mock.Anything).
		Once().
		Return(dnsResult, errors.New("dummy"))

	logger := &internalMock.Logger{}
	logger.On("WithField", "host", "wp.pl")
	logger.On("Debugf", mock.Anything, mock.Anything)

	checker := NewChecker(logger, dnsMock, host)
	result := checker.Check(context.TODO())
	assert.False(t, result.IsSuccess())
	assert.Equal(t, result.GetTime(), time.Duration(0))
	assert.True(t, dnsMock.AssertExpectations(t))
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
