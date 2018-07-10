package config

import (
	"time"

	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

// Config will hold already parsed information
type Config struct {
	RawConfig
	SuccessTimeThreshold  time.Duration
	SingleCheckTimeout    time.Duration
	MinimalCheckingPeriod time.Duration
	MaximalCheckingPeriod time.Duration
}

// DnsConfig holds all the information required to perform DNS checks
type DnsConfig struct {
	Hosts []string
}

type RawConfig struct {
	Dns                         DnsConfig
	Http                        HttpConfig `yaml:"http"`
	Persister                   Persister  `yaml:"persister"`
	SuccessRateThreshold        float32    `yaml:"success_rate_threshold"`
	SuccessTimeThresholdString  string     `yaml:"success_time_threshold"`
	SingleCheckTimeoutString    string     `yaml:"single_check_timeout"`
	MinimalCheckingPeriodString string     `yaml:"minimal_checking_period"`
	MaximalCheckingPeriodString string     `yaml:"maximal_checking_period"`
}

// HttpConfig is struct responsible to store all information about testing connection with HTTP
type HttpConfig struct {
	Urls []string `yaml:"urls"`
}

// Persister is a struct responsible to store all information about possible persisters (like DynamoDB)
type Persister struct {
	DynamoDB DynamoDbPersister `yaml:"dynamodb"`
}

// DynamoDbPersister is a struct responsible to store all information about AWS DynamoDB persister
type DynamoDbPersister struct {
	Enabled    bool   `yaml:"enabled"`
	Region     string `yaml:"region"`
	TableName  string `yaml:"table_name"`
	DeviceName string `yaml:"device_name"`
}

// Load will use 3rd party vendor to parse the file and return parsed config
func Load(fileName string) (Config, error) {
	rawConfig := RawConfig{}
	err := configor.Load(&rawConfig, fileName)
	if err != nil {
		return Config{}, errors.Wrapf(err, "Could not load config from file: '%s'", fileName)
	}

	return transformFromRawConfig(rawConfig)
}

// transformFromRawConfig will read the strings from RawConfig and parse them into Config struct
func transformFromRawConfig(rawConfig RawConfig) (Config, error) {
	c := Config{RawConfig: rawConfig}
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

	return c, nil
}
