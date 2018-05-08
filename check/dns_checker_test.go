// +build unit

package check

import (
	"testing"

	"time"

	"context"

	"github.com/krzysztof-gzocha/pingor/dns"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDNSChecker(t *testing.T) {
	successDnsMock := new(dnsMock)
	result := dns.DnsResult{Time: time.Second}
	successDnsMock.
		On("ResolveHost", "wp.pl").
		Once().
		Return(result, nil)

	checker := NewDNSChecker(successDnsMock, "wp.pl")
	assert.Equal(t, checker.dns, successDnsMock)
	assert.Len(t, checker.hosts, 1)
}

func TestDNSChecker_Check(t *testing.T) {
	successDnsMock := new(dnsMock)
	unsuccessfulDnsMock := new(dnsMock)
	dnsResult := dns.DnsResult{Time: time.Second}
	hosts := []string{"wp.pl", "onet.pl"}
	successDnsMock.
		On("ResolveHost", mock.Anything).
		Times(len(hosts)).
		Return(dnsResult, nil)
	unsuccessfulDnsMock.
		On("ResolveHost", mock.Anything).
		Times(len(hosts)).
		Return(dnsResult, errors.New("error"))

	checker := NewDNSChecker(successDnsMock, hosts...)
	result := checker.Check(context.TODO())
	assert.True(t, result.IsSuccess())
	assert.Equal(t, result.GetTime(), time.Second)
	assert.Len(t, result.GetSubResults(), 2)
	assert.True(t, successDnsMock.AssertExpectations(t))

	checker = NewDNSChecker(unsuccessfulDnsMock, hosts...)
	result = checker.Check(context.TODO())
	assert.False(t, result.IsSuccess())
	assert.Equal(t, result.GetTime(), time.Second)
	assert.Len(t, result.GetSubResults(), 2)
	assert.Equal(t, "Checking DNS with 2 hosts", result.GetMessage())
	assert.True(t, unsuccessfulDnsMock.AssertExpectations(t))

	checker = NewDNSChecker(successDnsMock)
	result = checker.Check(context.TODO())
	assert.False(t, result.IsSuccess())
	assert.Zero(t, result.GetSubResults())
	assert.Zero(t, result.GetSuccessRate())
	assert.Zero(t, result.GetTime())
}

type dnsMock struct {
	mock.Mock
	result dns.DnsResult
}

func (m dnsMock) ResolveHost(host string) (dns.DnsResult, error) {
	args := m.Called(host)

	return args.Get(0).(dns.DnsResult), args.Error(1)
}
