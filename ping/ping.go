package ping

import (
	"context"
	"net"
)

type PingInterface interface {
	Ping(ctx context.Context, ip net.IP) (Result, error)
}
