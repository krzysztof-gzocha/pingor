package main

import (
	"context"
	"flag"
	"net"
	"net/http"

	"github.com/Sirupsen/logrus"
	vendorsDynamoDb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/krzysztof-gzocha/pingor/pkg/check"
	"github.com/krzysztof-gzocha/pingor/pkg/check/dns"
	"github.com/krzysztof-gzocha/pingor/pkg/check/formatter/json"
	httpCheck "github.com/krzysztof-gzocha/pingor/pkg/check/http"
	"github.com/krzysztof-gzocha/pingor/pkg/check/multiple"
	"github.com/krzysztof-gzocha/pingor/pkg/check/periodic"
	"github.com/krzysztof-gzocha/pingor/pkg/config"
	"github.com/krzysztof-gzocha/pingor/pkg/event"
	"github.com/krzysztof-gzocha/pingor/pkg/log"
	"github.com/krzysztof-gzocha/pingor/pkg/metric"
	"github.com/krzysztof-gzocha/pingor/pkg/persister/aws/dynamodb"
	"github.com/krzysztof-gzocha/pingor/pkg/persister/aws/session"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber/reconnection"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// CLI Flags
	debugLevel := flag.Bool("debug", false, "Show debug level logs")
	configFile := flag.String("config", "config.yaml", "Specify config file name")
	flag.Parse()

	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.InfoLevel)
	if *debugLevel {
		logrusLogger.SetLevel(logrus.DebugLevel)
	}

	logger := log.NewLogrusWrapper(logrusLogger)

	// Configs
	cfg, err := config.Load(*configFile)
	if err != nil {
		logger.Errorf("Could not load config: %s", err.Error())
		return
	}

	if !cfg.Metrics.Enabled {
		run(context.Background(), cfg, logger)
		return
	}

	go run(context.Background(), cfg, logger)

	logger.Infof("Starting /metrics endpoint")
	http.Handle("/metrics", promhttp.Handler())
	httpErr := http.ListenAndServe(net.JoinHostPort("", cfg.Metrics.Port), nil)
	if httpErr != nil {
		logrus.Fatalf(httpErr.Error())
	}
}

func run(ctx context.Context, cfg config.Config, logger log.Logger) {
	// EventDispatcher with subscribers
	eventDispatcher := event.NewDispatcher()
	err := attachSubscribers(logger, eventDispatcher, cfg)
	if err != nil {
		logrus.Fatalf("Could not attach subscribers: %s", err.Error())
	}

	mainChecker := check.Checker(multiple.NewChecker(
		logger,
		cfg.SingleCheckTimeout,
		cfg.SuccessRateThreshold,
		cfg.SuccessTimeThreshold,
		getCheckers(logger, cfg)...,
	))

	if cfg.Metrics.Enabled {
		mainChecker = metric.NewInstrumentedSuccessRateChecker(mainChecker)
	}

	// Main checker
	checker := periodic.NewChecker(
		logger,
		eventDispatcher,
		mainChecker,
		cfg.MinimalCheckingPeriod,
		cfg.MaximalCheckingPeriod,
	)

	// Main logic
	checker.Check(ctx)
}

func attachSubscribers(logger log.Logger, dispatcher event.Dispatcher, cfg config.Config) error {
	checkLogger := subscriber.NewChecksLogger(logger)
	dispatcher.AttachSubscriber(periodic.ConnectionCheckEventName, checkLogger.LogConnectionCheckResult)

	reconnectSubscriber := subscriber.NewReconnection(dispatcher)
	dispatcher.AttachSubscriber(periodic.ConnectionCheckEventName, reconnectSubscriber.NotifyAboutReconnection)

	reconnectLogger := reconnection.NewLogger(logger, json.Formatter)
	dispatcher.AttachSubscriber(subscriber.ReconnectionEventName, reconnectLogger.LogReconnection)

	if !cfg.Persister.DynamoDB.Enabled {
		return nil
	}

	sess, err := session.CreateSession(cfg.Persister.DynamoDB.Region)
	if err != nil {
		return errors.Wrap(err, "Could not create AWS session")
	}

	persisterSubscriber := reconnection.NewPersister(
		logger,
		dynamodb.NewPersister(vendorsDynamoDb.New(sess), cfg.Persister.DynamoDB),
	)
	dispatcher.AttachSubscriber(subscriber.ReconnectionEventName, persisterSubscriber.PersistReconnectionEvent)

	return nil
}

func getCheckers(logger log.Logger, cfg config.Config) []check.Checker {
	checkers := make([]check.Checker, 0)

	if len(cfg.Dns.Hosts) > 0 {
		dnsClient := dns.Dns{}
		dnsCheckers := make([]check.Checker, 0)
		for _, url := range cfg.Dns.Hosts {
			var checker check.Checker = dns.NewChecker(logger, dnsClient, url)
			if cfg.Metrics.Enabled {
				checker = metric.NewInstrumentedDnsChecker(checker)
			}
			dnsCheckers = append(
				dnsCheckers,
				checker,
			)
		}

		checkers = append(checkers, multiple.NewChecker(
			logger,
			cfg.SingleCheckTimeout,
			cfg.SuccessRateThreshold,
			cfg.SuccessTimeThreshold,
			dnsCheckers...,
		))
	}

	if len(cfg.Http.Urls) > 0 {
		httpCheckers := make([]check.Checker, 0)
		for _, url := range cfg.Http.Urls {
			var checker check.Checker = httpCheck.NewChecker(logger, http.DefaultClient, url)
			if cfg.Metrics.Enabled {
				checker = metric.NewInstrumentedHttpChecker(checker)
			}

			httpCheckers = append(
				httpCheckers,
				checker,
			)
		}

		checkers = append(checkers, multiple.NewChecker(
			logger,
			cfg.SingleCheckTimeout,
			cfg.SuccessRateThreshold,
			cfg.SuccessTimeThreshold,
			httpCheckers...,
		))
	}

	return checkers
}
