// +build unit

package check

import (
	"context"
	"net"
	"testing"

	"time"

	"errors"

	"github.com/krzysztof-gzocha/pingor/ping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPingChecker(t *testing.T) {
	pingMock := pingMock{}
	ip := net.IPv4(1, 2, 3, 4)
	checker := NewPingChecker(pingMock, ip)
	assert.Equal(t, checker.ping, pingMock)
	assert.Equal(t, checker.ips[0], ip)
	assert.Len(t, checker.ips, 1)
}

func TestPingChecker_Check_Success(t *testing.T) {
	successPing := pingMock{}
	pingResult := ping.Result{Time: time.Second, PacketsReceived: 10, PacketsSent: 11}
	successPing.
		On("Ping", mock.Anything, mock.Anything).
		Once().
		Return(pingResult, nil)

	checker := NewPingChecker(successPing, net.IPv4(1, 1, 1, 1))
	result := checker.Check(context.TODO())

	assert.True(t, result.IsSuccess())
	assert.Equal(t, float32(10)/float32(11), result.GetSuccessRate())
	assert.Equal(t, time.Second, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 1)
}

func TestPingChecker_Check_Error(t *testing.T) {
	unsuccessfulPing := pingMock{}
	pingResult := ping.Result{Time: time.Second, PacketsReceived: 0, PacketsSent: 11}
	unsuccessfulPing.
		On("Ping", mock.Anything, mock.Anything).
		Once().
		Return(pingResult, errors.New("error"))

	checker := NewPingChecker(unsuccessfulPing, net.IPv4(1, 1, 1, 1))
	result := checker.Check(context.TODO())

	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0)/float32(11), result.GetSuccessRate())
	assert.Equal(t, time.Second, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 1)
}

func TestPingChecker_Check_EmptyIps(t *testing.T) {
	unsuccessfulPing := pingMock{}
	pingResult := ping.Result{Time: time.Second, PacketsReceived: 0, PacketsSent: 11}
	unsuccessfulPing.
		On("Ping", mock.Anything, mock.Anything).
		Times(0).
		Return(pingResult, errors.New("error"))

	checker := NewPingChecker(unsuccessfulPing)
	result := checker.Check(context.TODO())

	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0), result.GetSuccessRate())
	assert.Equal(t, time.Duration(0), result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 0)
}

type pingMock struct {
	mock.Mock
}

func (m pingMock) Ping(ctx context.Context, ip net.IP) (ping.Result, error) {
	args := m.Called(ctx, ip)

	return args.Get(0).(ping.Result), args.Error(1)
}
