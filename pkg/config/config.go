package config

import (
	"net"
	"time"

	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

// Config will hold already parsed information
type Config struct {
	rawConfig
	Ping                  PingConfig
	SuccessTimeThreshold  time.Duration
	SingleCheckTimeout    time.Duration
	MinimalCheckingPeriod time.Duration
	MaximalCheckingPeriod time.Duration
}

// DnsConfig holds all the information required to perform DNS checks
type DnsConfig struct {
	Hosts []string
}

// PingConfig holds all the information required to perform ping checks
type PingConfig struct {
	IPs []net.IP
}

type rawConfig struct {
	Dns                         DnsConfig
	RawPing                     rawPingConfig `yaml:"ping"`
	Http                        HttpConfig    `yaml:"http"`
	Persister                   Persister     `yaml:"persister"`
	SuccessRateThreshold        float32       `yaml:"success_rate_threshold"`
	SuccessTimeThresholdString  string        `yaml:"success_time_threshold"`
	SingleCheckTimeoutString    string        `yaml:"single_check_timeout"`
	MinimalCheckingPeriodString string        `yaml:"minimal_checking_period"`
	MaximalCheckingPeriodString string        `yaml:"maximal_checking_period"`
}

type rawPingConfig struct {
	IpStrings []string `yaml:"ips"`
}

type HttpConfig struct {
	Urls []string `yaml:"urls"`
}

type Persister struct {
	DynamoDB DynamoDbPersister `yaml:"dynamodb"`
}

type DynamoDbPersister struct {
	Enabled    bool   `yaml:"enabled"`
	Region     string `yaml:"region"`
	TableName  string `yaml:"table_name"`
	DeviceName string `yaml:"device_name"`
}

// Load will use 3rd party vendor to parse the file and return parsed config
func Load(fileName string) (Config, error) {
	rawConfig := rawConfig{}
	err := configor.Load(&rawConfig, fileName)
	if err != nil {
		return Config{}, errors.Wrapf(err, "Could not load config from file: '%s'", fileName)
	}

	return transformFromRawConfig(rawConfig)
}

// transformFromRawConfig will read the strings from RawConfig and parse them into Config struct
func transformFromRawConfig(rawConfig rawConfig) (Config, error) {
	c := Config{rawConfig: rawConfig}
	successTime, err := time.ParseDuration(rawConfig.SuccessTimeThresholdString)
	if err != nil {
		return c, errors.Wrapf(err, "Could not parse success time threshold: %s", rawConfig.SuccessTimeThresholdString)
	}
	c.SuccessTimeThreshold = successTime

	singleCheckTimeout, err := time.ParseDuration(rawConfig.SingleCheckTimeoutString)
	if err != nil {
		return c, errors.Wrapf(err, "Could not parse single check timeout: %s", rawConfig.SuccessTimeThresholdString)
	}
	c.SingleCheckTimeout = singleCheckTimeout

	minimalCheckingPeriod, err := time.ParseDuration(rawConfig.MinimalCheckingPeriodString)
	if err != nil {
		return c, errors.Wrapf(err, "Could not parse minimal checking period: %s", rawConfig.SuccessTimeThresholdString)
	}
	c.MinimalCheckingPeriod = minimalCheckingPeriod

	maximalCheckingPeriod, err := time.ParseDuration(rawConfig.MaximalCheckingPeriodString)
	if err != nil {
		return c, errors.Wrapf(err, "Could not parse maximal checking period: %s", rawConfig.SuccessTimeThresholdString)
	}
	c.MaximalCheckingPeriod = maximalCheckingPeriod
	c.Ping = transformFromRawPingConfig(c.RawPing)

	return c, nil
}

// transformFromRawPingConfig will read the strings from RawPingConfig and parse them into PingConfig struct
func transformFromRawPingConfig(rpc rawPingConfig) PingConfig {
	results := make([]net.IP, len(rpc.IpStrings))
	for k, ipString := range rpc.IpStrings {
		results[k] = net.ParseIP(ipString)
	}

	return PingConfig{IPs: results}
}
