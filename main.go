package main

import (
	"context"
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/check"
	"github.com/krzysztof-gzocha/pingor/config"
	"github.com/krzysztof-gzocha/pingor/dns"
	"github.com/krzysztof-gzocha/pingor/event"
	"github.com/krzysztof-gzocha/pingor/ping"
	"github.com/krzysztof-gzocha/pingor/subscriber"
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
	eventDispatcher.AttachSubscriber(check.ConnectionCheckEventName, subscriber.LogConnectionCheckResult)
	reconnectSubscriber := subscriber.NewReconnectionSubscriber(check.JsonResultPrinter)
	eventDispatcher.AttachSubscriber(check.ConnectionCheckEventName, reconnectSubscriber.NotifyAboutReconnection)

	// Main checker
	checker := check.NewPeriodicCheckerWrapper(
		eventDispatcher,
		check.NewMultipleChecker(
			cfg.SingleCheckTimeout,
			check.NewPingChecker(ping.PingCommand{}, cfg.Ping.IPs...),
			check.NewDNSChecker(dns.Dns{}, cfg.Dns.Hosts...),
		),
		cfg.SuccessRateThreshold,
		cfg.SuccessTimeThreshold,
		cfg.MinimalCheckingPeriod,
		cfg.MaximalCheckingPeriod,
	)

	// Main logic
	checker.Check(context.Background())
}

func configureLogLevel(verboseLevel bool) {
	logrus.SetLevel(logrus.InfoLevel)
	if verboseLevel {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
