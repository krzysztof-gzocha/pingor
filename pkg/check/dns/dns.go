package dns

import (
	"net"
	"time"
)

// Resolver should be implemented by any service capable of resolving host into IP
type Resolver interface {
	ResolveHost(host string) (Result, error)
}

// Result is data-transfer struct that holds the result of resolving host name into IP and the time it took
type Result struct {
	Ip   *net.IPAddr
	Time time.Duration
}

// Dns is struct capable to resolve host to IP
type Dns struct{}

// ResolveHost will try to resolve the IP from provided host name
func (d Dns) ResolveHost(host string) (Result, error) {
	start := time.Now()
	ip, err := net.ResolveIPAddr("ip4:icmp", host)
	end := time.Now()

	return Result{
		Ip:   ip,
		Time: end.Sub(start),
	}, err
}
