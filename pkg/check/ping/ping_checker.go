package ping

import (
	"context"
	"net"

	"fmt"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
)

// Checker will run 'ping' command on underlying system to check internet connection and interpret it's response
type Checker struct {
	ping PingerInterface
	ips  []net.IP
}

// NewChecker will return new instance of Checker
func NewChecker(ping PingerInterface, ips ...net.IP) Checker {
	return Checker{ping: ping, ips: ips}
}

// Check will run 'ping' command on underlying system to check internet connection and interpret it's response
// Result's time is average time of all the tests.
func (p Checker) Check(ctx context.Context) result.ResultInterface {
	if len(p.ips) == 0 {
		return result.Result{Success: false, Message: fmt.Sprintf("Checking ping command with %d IPs", len(p.ips))}
	}

	overallResult := result.Result{Success: true, Message: fmt.Sprintf("Checking ping command with %d IPs", len(p.ips))}
	for _, ip := range p.ips {
		result := p.singleCheck(ctx, ip)
		overallResult.SubResults = append(overallResult.SubResults, result)
	}

	return p.calculateOverallResult(overallResult)
}

func (p Checker) singleCheck(ctx context.Context, ip net.IP) result.Result {
	result := result.Result{Success: true}
	logrus.Debugf("Checker: starting to check %s", ip.String())
	pingResult, err := p.ping.Ping(ctx, ip)
	result.Message = fmt.Sprintf("%T:%s", p, ip.String())
	if err != nil {
		errMsg := fmt.Sprintf("%T:%s: %s", p, ip.String(), err.Error())
		result.Message = errMsg
		result.Success = false
	}

	if !pingResult.AtLeastOneSuccess() {
		result.Success = false
	}

	if pingResult.PacketsReceived > 0 {
		result.SuccessRate = float32(pingResult.PacketsReceived) / float32(pingResult.PacketsSent)
	}

	result.Time = pingResult.Time

	return result
}

func (p Checker) calculateOverallResult(overallResult result.Result) result.Result {
	var successRates float32
	for _, subResult := range overallResult.SubResults {
		successRates += subResult.GetSuccessRate()
		if !subResult.IsSuccess() {
			overallResult.Success = false
		}
		overallResult.Time += subResult.GetTime()
	}
	overallResult.SuccessRate = successRates / float32(len(overallResult.SubResults))
	overallResult.Time = overallResult.Time / time.Duration(len(overallResult.SubResults))

	return overallResult
}
