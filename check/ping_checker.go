package check

import (
	"context"
	"net"

	"fmt"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/ping"
)

// PingChecker will run 'ping' command on underlying system to check internet connection and interpret it's response
type PingChecker struct {
	ping ping.PingInterface
	ips  []net.IP
}

// NewPingChecker will return new instance of PingChecker
func NewPingChecker(ping ping.PingInterface, ips ...net.IP) PingChecker {
	return PingChecker{ping: ping, ips: ips}
}

// Check will run 'ping' command on underlying system to check internet connection and interpret it's response
// Result's time is average time of all the tests.
func (p PingChecker) Check(ctx context.Context) ResultInterface {
	if len(p.ips) == 0 {
		return Result{Success: false, Message: fmt.Sprintf("Checking ping command with %d IPs", len(p.ips))}
	}

	overallResult := Result{Success: true, Message: fmt.Sprintf("Checking ping command with %d IPs", len(p.ips))}
	for _, ip := range p.ips {
		result := Result{Success: true}
		logrus.Debugf("PingChecker: starting to check %s", ip.String())
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
		overallResult.SubResults = append(overallResult.SubResults, result)
	}

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
