package ping

import (
	"context"
	"net"
)

// PingInterface should be implemented by any service capable to send and interpret ping
type PingInterface interface {
	Ping(ctx context.Context, ip net.IP) (Result, error)
}
