package dns

import (
	"net"
	"time"
)

// DnsInterface should be implemented by any service capable of resolving host into IP
type DnsInterface interface {
	ResolveHost(host string) (DnsResult, error)
}

// Dns is struct capable to resolve host to IP
type Dns struct{}

// DnsResult is data-transfer struct that holds the result of resolving host name into IP and the time it took
type DnsResult struct {
	Ip   *net.IPAddr
	Time time.Duration
}

// ResolveHost will try to resolve the IP from provided host name
func (d Dns) ResolveHost(host string) (DnsResult, error) {
	start := time.Now()
	ip, err := net.ResolveIPAddr("ip4:icmp", host)
	end := time.Now()

	return DnsResult{
		Ip:   ip,
		Time: end.Sub(start),
	}, err
}
