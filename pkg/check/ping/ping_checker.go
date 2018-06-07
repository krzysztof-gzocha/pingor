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
		res := p.singleCheck(ctx, ip)
		overallResult.SubResults = append(overallResult.SubResults, res)
	}

	p.calculateOverallResult(&overallResult)
	logrus.WithField("successRate", overallResult.SuccessRate*100).Debugf("%T: done", p)

	return overallResult
}

func (p Checker) singleCheck(ctx context.Context, ip net.IP) result.Result {
	res := result.Result{Success: true}
	logrus.WithField("ip", ip.String()).Debugf("%T: starting to check", p)
	pingResult, err := p.ping.Ping(ctx, ip)
	res.Message = fmt.Sprintf("%T:%s", p, ip.String())
	if err != nil {
		errMsg := fmt.Sprintf("%T:%s: %s", p, ip.String(), err.Error())
		res.Message = errMsg
		res.Success = false
	}

	if !pingResult.AtLeastOneSuccess() {
		res.Success = false
	}

	if pingResult.PacketsReceived > 0 {
		res.SuccessRate = float32(pingResult.PacketsReceived) / float32(pingResult.PacketsSent)
	}

	res.Time = pingResult.Time

	return res
}

func (p Checker) calculateOverallResult(overallResult *result.Result) {
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
}
