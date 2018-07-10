// +build unit

package main

import (
	"context"
	"testing"
	"time"

	"github.com/krzysztof-gzocha/pingor/pkg/config"
	"github.com/krzysztof-gzocha/pingor/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestRun_NotPanics_WithDynamoDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
	defer cancel()

	cfg := config.Config{
		RawConfig: config.RawConfig{
			Dns:  config.DnsConfig{Hosts: []string{"wp.pl"}},
			Http: config.HttpConfig{Urls: []string{"wp.pl"}},
			Persister: config.Persister{DynamoDB: config.DynamoDbPersister{
				Enabled: true,
				Region:  "test",
			}},
		},
	}

	assert.NotPanics(t, func() {
		run(ctx, cfg, &log.Nil{})
	})
}

func TestRun_NotPanics_WithoutDynamoDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
	defer cancel()

	cfg := config.Config{
		RawConfig: config.RawConfig{
			Dns:  config.DnsConfig{Hosts: []string{"wp.pl"}},
			Http: config.HttpConfig{Urls: []string{"wp.pl"}},
		},
	}

	assert.NotPanics(t, func() {
		run(ctx, cfg, &log.Nil{})
	})
}
