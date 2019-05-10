package metric

import (
	"context"

	"github.com/krzysztof-gzocha/pingor/pkg/check"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpSuccessRate = promauto.NewGauge(
		prometheus.GaugeOpts(prometheus.Opts{
			Name: "http_pingor_success_rate",
			Help: "HTTP response success rate in percents",
		}),
	)

	httpResponsesTimes = promauto.NewGaugeVec(prometheus.GaugeOpts(prometheus.Opts{
		Name: "http_pingor_response_times_ns",
		Help: "HTTP response time presented in nano seconds",
	}), []string{"url"})

	dnsResponsesTimes = promauto.NewGaugeVec(prometheus.GaugeOpts(prometheus.Opts{
		Name: "dns_pingor_response_times_ns",
		Help: "DNS response time presented in nano seconds",
	}), []string{"url"})
)

type InstrumentedGaugeChecker struct {
	checker check.Checker
	gauge   *prometheus.GaugeVec
}

type InstrumentedSuccessRateChecker struct {
	checker check.Checker
}

func NewInstrumentedHttpChecker(checker check.Checker) *InstrumentedGaugeChecker {
	return &InstrumentedGaugeChecker{
		checker: checker,
		gauge:   httpResponsesTimes,
	}
}

func NewInstrumentedDnsChecker(checker check.Checker) *InstrumentedGaugeChecker {
	return &InstrumentedGaugeChecker{
		checker: checker,
		gauge:   dnsResponsesTimes,
	}
}

func NewInstrumentedSuccessRateChecker(checker check.Checker) *InstrumentedSuccessRateChecker {
	return &InstrumentedSuccessRateChecker{
		checker: checker,
	}
}

func (i *InstrumentedGaugeChecker) Check(ctx context.Context) result.Result {
	checkResult := i.checker.Check(ctx)
	i.gauge.
		WithLabelValues(checkResult.GetURL()).
		Set(float64(checkResult.GetTime().Nanoseconds()))

	return checkResult
}

func (i *InstrumentedSuccessRateChecker) Check(ctx context.Context) result.Result {
	checkResult := i.checker.Check(ctx)
	httpSuccessRate.Set(float64(checkResult.GetSuccessRate()))

	return checkResult
}
