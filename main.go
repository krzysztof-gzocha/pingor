package main

import (
	"context"
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/pkg/check"
	"github.com/krzysztof-gzocha/pingor/pkg/check/dns"
	"github.com/krzysztof-gzocha/pingor/pkg/check/multiple"
	"github.com/krzysztof-gzocha/pingor/pkg/check/periodic"
	"github.com/krzysztof-gzocha/pingor/pkg/check/ping"
	"github.com/krzysztof-gzocha/pingor/pkg/check/printer/json"
	"github.com/krzysztof-gzocha/pingor/pkg/config"
	"github.com/krzysztof-gzocha/pingor/pkg/event"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
)

func main() {
	// CLI Flags
	debugLevel := flag.Bool("debug", false, "Show debug level logs")
	configFile := flag.String("config", "config.yaml", "Specify config file name")
	flag.Parse()

	// Configs
	configureLogLevel(*debugLevel)
	cfg, err := config.Load(*configFile)
	if err != nil {
		logrus.Fatalf("Could not load config: %s", err.Error())
	}

	// EventDispatcher with subscribers
	eventDispatcher := event.NewDispatcher()
	eventDispatcher.AttachSubscriber(periodic.ConnectionCheckEventName, subscriber.LogConnectionCheckResult)
	reconnectSubscriber := subscriber.NewReconnectionSubscriber(json.Printer)
	eventDispatcher.AttachSubscriber(periodic.ConnectionCheckEventName, reconnectSubscriber.NotifyAboutReconnection)

	// Main checker
	checker := periodic.NewChecker(
		eventDispatcher,
		multiple.NewChecker(
			cfg.SingleCheckTimeout,
			cfg.SuccessRateThreshold,
			cfg.SuccessTimeThreshold,
			getCheckers(cfg)...,
		),
		cfg.MinimalCheckingPeriod,
		cfg.MaximalCheckingPeriod,
	)

	// Main logic
	checker.Check(context.Background())
}

func getCheckers(cfg config.Config) []check.CheckerInterface {
	checkers := make([]check.CheckerInterface, 0)

	if len(cfg.Ping.IPs) > 0 {
		checkers = append(checkers, ping.NewChecker(ping.Command{}, cfg.Ping.IPs...))
	}

	if len(cfg.Dns.Hosts) > 0 {
		checkers = append(checkers, dns.NewChecker(dns.Dns{}, cfg.Dns.Hosts...))
	}

	return checkers
}

func configureLogLevel(verboseLevel bool) {
	logrus.SetLevel(logrus.InfoLevel)
	if verboseLevel {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
