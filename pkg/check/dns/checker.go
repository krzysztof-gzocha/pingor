package dns

import (
	"context"
	"fmt"
	"time"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/log"
)

// Checker will try to resolve provided hosts into IPs in order to check the connection to DNS
type Checker struct {
	logger log.LoggerInterface
	dns    ResolverInterface
	hosts  []string
}

// NewChecker will return new instance of Checker
func NewChecker(logger log.LoggerInterface, dns ResolverInterface, hosts ...string) Checker {
	return Checker{logger: logger, dns: dns, hosts: hosts}
}

// Check will try to resolve provided hosts into IPs in order to check the connection to DNS.
// Time result is average time required to resolve all the hosts.
func (c Checker) Check(ctx context.Context) result.ResultInterface {
	if len(c.hosts) == 0 {
		return result.Result{}
	}

	overallResult := result.Result{Success: true, Message: fmt.Sprintf("Checking DNS with %d hosts", len(c.hosts))}
	for _, host := range c.hosts {
		overallResult.SubResults = append(overallResult.SubResults, c.singleCheck(ctx, host))
	}

	overallResult = c.calculateOverallChecker(overallResult)
	c.logger.
		WithField("successRate", overallResult.SuccessRate*100).
		Debugf("%T: done", c)

	return overallResult
}

func (c Checker) singleCheck(ctx context.Context, host string) result.ResultInterface {
	res := result.Result{Success: true}
	c.logger.
		WithField("host", host).
		Debugf("%T: starting to check", c)

	dnsResult, err := c.dns.ResolveHost(host)
	res.Time = dnsResult.Time
	res.Message = fmt.Sprintf("%T:%s", c, host)
	res.SuccessRate = 1
	if err != nil {
		errMsg := fmt.Sprintf("%T:%s: Failed to resolve host: %s", c, host, err.Error())
		res.Success = false
		res.Message = errMsg
		res.SuccessRate = 0
	}

	return res
}

func (c Checker) calculateOverallChecker(overallResult result.Result) result.Result {
	var successResults float32
	var totalTime time.Duration
	for _, subresult := range overallResult.SubResults {
		if subresult.IsSuccess() {
			successResults++
		} else {
			overallResult.Success = false
		}
		totalTime += subresult.GetTime()
	}

	overallResult.SuccessRate = successResults / float32(len(c.hosts))
	overallResult.Time = totalTime / time.Duration(len(c.hosts))

	return overallResult
}
