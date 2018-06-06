package dns

import (
	"context"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
)

// Checker will try to resolve provided hosts into IPs in order to check the connection to DNS
type Checker struct {
	dns   ResolverInterface
	hosts []string
}

// NewChecker will return new instance of Checker
func NewChecker(dns ResolverInterface, hosts ...string) Checker {
	return Checker{dns: dns, hosts: hosts}
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
	logrus.Debugf("Checker: success rate: %.2f", overallResult.SuccessRate*100)

	return overallResult
}

func (c Checker) singleCheck(ctx context.Context, host string) result.Result {
	result := result.Result{Success: true}
	logrus.Debugf("Checker: starting to check: %s", host)

	dnsResult, err := c.dns.ResolveHost(host)
	result.Time = dnsResult.Time
	result.Message = fmt.Sprintf("%T:%s", c, host)
	result.SuccessRate = 1
	if err != nil {
		errMsg := fmt.Sprintf("%T:%s: Failed to resolve host: %s", c, host, err.Error())
		result.Success = false
		result.Message = errMsg
		result.SuccessRate = 0
	}

	return result
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
