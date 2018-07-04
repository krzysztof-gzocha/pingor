// +build unit

package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cfg, err := Load("../../config.yaml")
	assert.Nil(t, err)
	assert.NotZero(t, cfg.SuccessRateThreshold)

	cfg, err = Load("../../main.go")
	assert.Error(t, err)
	assert.Zero(t, cfg.SuccessRateThreshold)
}

func TestTransformation(t *testing.T) {
	rawConfig := rawConfig{
		SuccessTimeThresholdString:  "10s",
		SingleCheckTimeoutString:    "10m",
		MinimalCheckingPeriodString: "100ms",
		MaximalCheckingPeriodString: "1h",
	}

	config, err := transformFromRawConfig(rawConfig)
	assert.Nil(t, err)
	assert.Equal(t, config.SingleCheckTimeout, time.Minute*10)
	assert.Equal(t, config.SuccessTimeThreshold, time.Second*10)
	assert.Equal(t, config.MinimalCheckingPeriod, time.Millisecond*100)
	assert.Equal(t, config.MaximalCheckingPeriod, time.Hour)
}

func TestTransformationErrors(t *testing.T) {
	correctDuration := "1m"
	incorrectDuration := "some bad string"

	configs := []rawConfig{
		{SuccessTimeThresholdString: incorrectDuration},
		{
			SuccessTimeThresholdString: correctDuration,
			SingleCheckTimeoutString:   incorrectDuration,
		},
		{
			SuccessTimeThresholdString:  correctDuration,
			SingleCheckTimeoutString:    correctDuration,
			MinimalCheckingPeriodString: incorrectDuration,
		},
		{
			SuccessTimeThresholdString:  correctDuration,
			SingleCheckTimeoutString:    correctDuration,
			MinimalCheckingPeriodString: correctDuration,
			MaximalCheckingPeriodString: incorrectDuration,
		},
	}

	for _, config := range configs {
		_, err := transformFromRawConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Could not parse")
	}
}
