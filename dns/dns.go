package dns

import (
	"net"
	"time"
)

type DnsInterface interface {
	ResolveHost(host string) (DnsResult, error)
}

type Dns struct{}

type DnsResult struct {
	Ip   *net.IPAddr
	Time time.Duration
}

func (d Dns) ResolveHost(host string) (DnsResult, error) {
	start := time.Now()
	ip, err := net.ResolveIPAddr("ip4:icmp", host)
	end := time.Now()

	return DnsResult{
		Ip:   ip,
		Time: end.Sub(start),
	}, err
}
