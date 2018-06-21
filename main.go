package main

import (
	"context"
	"flag"

	"net/http"

	"github.com/Sirupsen/logrus"
	vendorsDynamoDb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/krzysztof-gzocha/pingor/pkg/check"
	"github.com/krzysztof-gzocha/pingor/pkg/check/dns"
	"github.com/krzysztof-gzocha/pingor/pkg/check/formatter/json"
	httpCheck "github.com/krzysztof-gzocha/pingor/pkg/check/http"
	"github.com/krzysztof-gzocha/pingor/pkg/check/multiple"
	"github.com/krzysztof-gzocha/pingor/pkg/check/periodic"
	"github.com/krzysztof-gzocha/pingor/pkg/check/ping"
	"github.com/krzysztof-gzocha/pingor/pkg/config"
	"github.com/krzysztof-gzocha/pingor/pkg/event"
	"github.com/krzysztof-gzocha/pingor/pkg/persister/aws/dynamodb"
	"github.com/krzysztof-gzocha/pingor/pkg/persister/aws/session"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber/reconnection"
	"github.com/pkg/errors"
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
	err = attachSubscribers(eventDispatcher, cfg)
	if err != nil {
		logrus.Fatalf("Could not attach subscribers: %s", err.Error())
	}

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

func attachSubscribers(dispatcher event.DispatcherInterface, cfg config.Config) error {
	dispatcher.AttachSubscriber(periodic.ConnectionCheckEventName, subscriber.LogConnectionCheckResult)

	reconnectSubscriber := subscriber.NewReconnection(dispatcher)
	dispatcher.AttachSubscriber(periodic.ConnectionCheckEventName, reconnectSubscriber.NotifyAboutReconnection)

	reconnectLogger := reconnection.NewLogger(json.Formatter)
	dispatcher.AttachSubscriber(subscriber.ReconnectionEventName, reconnectLogger.LogReconnection)

	if !cfg.Persister.DynamoDB.Enabled {
		return nil
	}

	sess, err := session.CreateSession()
	if err != nil {
		return errors.Wrap(err, "Could not create AWS session")
	}

	persisterSubscriber := reconnection.NewPersister(
		dynamodb.NewPersister(vendorsDynamoDb.New(sess), cfg.Persister.DynamoDB),
	)
	dispatcher.AttachSubscriber(subscriber.ReconnectionEventName, persisterSubscriber.PersistReconnectionEvent)

	return nil
}

func getCheckers(cfg config.Config) []check.CheckerInterface {
	checkers := make([]check.CheckerInterface, 0)

	if len(cfg.Ping.IPs) > 0 {
		checkers = append(checkers, ping.NewChecker(ping.Command{}, cfg.Ping.IPs...))
	}

	if len(cfg.Dns.Hosts) > 0 {
		checkers = append(checkers, dns.NewChecker(dns.Dns{}, cfg.Dns.Hosts...))
	}

	if len(cfg.Http.Urls) > 0 {
		checkers = append(checkers, httpCheck.NewChecker(http.DefaultClient, cfg.Http.Urls...))
	}

	return checkers
}

func configureLogLevel(verboseLevel bool) {
	logrus.SetLevel(logrus.InfoLevel)
	if verboseLevel {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
