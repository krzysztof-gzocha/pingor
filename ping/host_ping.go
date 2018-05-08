package ping

import (
	"context"
	"net"
	"os/exec"

	"regexp"

	"strconv"

	"time"

	"github.com/pkg/errors"
)

// Result is single ping result data transfer object
type Result struct {
	PacketsSent     uint
	PacketsReceived uint
	Time            time.Duration
	IP              net.IP
}

// AtLeastOneSuccess will return true if success rate if higher than 0
func (r Result) AtLeastOneSuccess() bool {
	return r.SuccessRate() > 0
}

// SuccessRate will return ratio of packets sent to packets received as float32 >=0 and <=1
func (r Result) SuccessRate() float32 {
	if r.PacketsReceived == 0 {
		return 0
	}

	return float32(r.PacketsSent) / float32(r.PacketsReceived)
}

// PingCommand is service that will call ping command on the host and interpret it's response
type PingCommand struct{}

// Ping will run ping command on the host OS.
func (p PingCommand) Ping(ctx context.Context, ip net.IP) (Result, error) {
	cmd := exec.CommandContext(ctx, "ping", "-q", "-c", "3", ip.String())
	byteOutput, err := cmd.Output()
	if err != nil {
		return Result{}, errors.Wrap(err, "Error occurred while reading the output of ping command")
	}

	result, err := p.parseOutput(byteOutput)
	result.IP = ip
	if err != nil {
		return result, err
	}

	return result, nil
}

func (p PingCommand) parseOutput(output []byte) (Result, error) {
	stringOutput := string(output)
	reg, err := regexp.Compile(`(\d{1,4}) packets transmitted, (\d{1,4}) received`)
	if err != nil {
		return Result{}, err
	}

	found := reg.FindStringSubmatch(stringOutput)
	if len(found) != 3 {
		return Result{}, errors.Errorf("Didn't understand output for parsing number of packets: %s", stringOutput)
	}

	transmitted, err := strconv.Atoi(found[1])
	if err != nil {
		return Result{}, errors.Wrapf(err, "Could not convert packets transmitted '%s' to int", found[1])
	}

	received, err := strconv.Atoi(found[2])
	if err != nil {
		return Result{}, errors.Wrapf(err, "Could not convert packets received '%s' to int", found[2])
	}

	timeReg, err := regexp.Compile(`rtt min/avg/max/mdev = [\d\.]{1,6}/([\d\.]{1,6})/[\d\.]{1,6}/[\d\.]{1,6} ms`)
	if err != nil {
		return Result{}, err
	}

	found = timeReg.FindStringSubmatch(stringOutput)
	if len(found) != 2 {
		return Result{}, errors.Errorf("Didn't understand output for parsing time: %s", stringOutput)
	}

	parsedTime, err := strconv.ParseFloat(found[1], 32)
	if err != nil {
		return Result{}, errors.Wrapf(err, "Could not convert packets transmitted '%s' to int", found[1])
	}

	return Result{
		PacketsSent:     uint(transmitted),
		PacketsReceived: uint(received),
		Time:            time.Microsecond * time.Duration(parsedTime*1000),
	}, nil
}
