// +build unit

package config

import (
	"net"
	"testing"
	"time"
)

func TestTransformation(t *testing.T) {
	rawConfig := rawConfig{
		SuccessTimeThresholdString:  "10s",
		SingleCheckTimeoutString:    "10m",
		MinimalCheckingPeriodString: "100ms",
		MaximalCheckingPeriodString: "1h",
		RawPing: rawPingConfig{
			IpStrings: []string{"1.1.1.1", "192.168.100.10"},
		},
	}

	config, err := transformFromRawConfig(rawConfig)
	if err != nil {
		t.Fatalf("Transformation failed: %s", err.Error())
	}

	if config.SingleCheckTimeout != time.Minute*10 {
		t.Fatalf("SingleCheckTimeout was parsed incorectly")
	}

	if config.SuccessTimeThreshold != time.Second*10 {
		t.Fatalf("SuccessTimeThreshold was parsed incorectly")
	}

	if config.MinimalCheckingPeriod != time.Millisecond*100 {
		t.Fatalf("MinimalCheckingPeriod was parsed incorectly")
	}

	if config.MaximalCheckingPeriod != time.Hour {
		t.Fatalf("MaximalCheckingPeriod was parsed incorectly")
	}

	if len(config.Ping.IPs) != len(rawConfig.RawPing.IpStrings) {
		t.Fatalf("Transformed ping IPs: %d, raw IPs: %d", len(config.Ping.IPs), len(rawConfig.RawPing.IpStrings))
	}
}

func TestPingTransformation(t *testing.T) {
	rawPingConfig := rawPingConfig{
		IpStrings: []string{"1.1.1.1", "182.123.231.23"},
	}

	pingConfig := transformFromRawPingConfig(rawPingConfig)

	if len(pingConfig.IPs) != len(rawPingConfig.IpStrings) {
		t.Fatalf("Transformed ping IPs: %d, raw IPs: %d", len(pingConfig.IPs), len(rawPingConfig.IpStrings))
	}

	if !pingConfig.IPs[0].Equal(net.IPv4(1, 1, 1, 1)) {
		t.Fatalf("First IP is badly parsed: %s", pingConfig.IPs[0].String())
	}

	if !pingConfig.IPs[1].Equal(net.IPv4(182, 123, 231, 23)) {
		t.Fatalf("Second IP is badly parsed: %s", pingConfig.IPs[1].String())
	}
}
