package ping

import (
	"context"
	"net"
)

// PingerInterface should be implemented by any service capable to send and interpret ping
type PingerInterface interface {
	Ping(ctx context.Context, ip net.IP) (Result, error)
}
