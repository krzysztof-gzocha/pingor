package check

import (
	"context"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/dns"
)

// DNSChecker will try to resolve provided hosts into IPs in order to check the connection to DNS
type DNSChecker struct {
	dns   dns.DnsInterface
	hosts []string
}

// NewDNSChecker will return new instance of DNSChecker
func NewDNSChecker(dns dns.DnsInterface, hosts ...string) DNSChecker {
	return DNSChecker{dns: dns, hosts: hosts}
}

// Check will try to resolve provided hosts into IPs in order to check the connection to DNS.
// Time result is average time required to resolve all the hosts.
func (c DNSChecker) Check(ctx context.Context) ResultInterface {
	if len(c.hosts) == 0 {
		return Result{}
	}

	overallResult := Result{Success: true, Message: fmt.Sprintf("Checking DNS with %d hosts", len(c.hosts))}
	for _, host := range c.hosts {
		result := Result{Success: true}
		logrus.Debugf("DNSChecker: starting to check: %s", host)

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

		overallResult.SubResults = append(overallResult.SubResults, result)
	}

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
	logrus.Debugf("DNSChecker: success rate: %.2f", overallResult.SuccessRate*100)

	return overallResult
}
