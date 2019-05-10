package dns

import (
	"context"
	"fmt"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/log"
)

// Checker will try to resolve provided hosts into IPs in order to check the connection to DNS
type Checker struct {
	logger log.Logger
	dns    Resolver
	host   string
}

// NewChecker will return new instance of Checker
func NewChecker(logger log.Logger, dns Resolver, host string) Checker {
	return Checker{logger: logger, dns: dns, host: host}
}

// Check will try to resolve provided hosts into IPs in order to check the connection to DNS.
// Time result is average time required to resolve all the hosts.
func (c Checker) Check(ctx context.Context) result.Result {
	res := result.DefaultResult{URL: c.host, Success: true}
	c.logger.
		WithField("host", c.host).
		Debugf("%T: starting to check", c)

	dnsResult, err := c.dns.ResolveHost(c.host)
	if err != nil {
		errMsg := fmt.Sprintf("%T:%s: Failed to resolve host: %s", c, c.host, err.Error())
		res.Success = false
		res.Message = errMsg
		res.SuccessRate = 0

		return res
	}

	res.Time = dnsResult.Time
	res.Message = fmt.Sprintf("%T:%s", c, c.host)
	res.SuccessRate = 1

	return res
}
