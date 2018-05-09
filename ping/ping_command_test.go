// +build unit

package ping

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingCommand_Ping(t *testing.T) {
	pc := PingCommand{}
	assert.Implements(t, (*PingInterface)(nil), pc)
	result, err := pc.Ping(context.TODO(), net.ParseIP("8.8.8.8"))
	assert.NoError(t, err)
	assert.NotZero(t, result.Time)
	assert.True(t, result.PacketsSent > 0)
	assert.True(t, result.PacketsReceived > 0)
	assert.True(t, result.AtLeastOneSuccess())
	assert.True(t, result.SuccessRate() > 0)
	assert.NotEmpty(t, result.IP)
}

func TestPingCommand_Ping_Error(t *testing.T) {
	pc := PingCommand{}
	result, err := pc.Ping(context.TODO(), net.ParseIP("something"))
	assert.Error(t, err)
	assert.Zero(t, result.Time)
	assert.Zero(t, result.PacketsSent)
	assert.Zero(t, result.PacketsReceived)
	assert.Empty(t, result.IP)

	assert.Equal(t, float32(0), result.SuccessRate())
}
