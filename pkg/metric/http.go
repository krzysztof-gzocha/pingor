package metric

import (
	"github.com/krzysztof-gzocha/pingor/pkg/check/http"
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

	// dnsResponsesTimes = promauto.NewGaugeVec(prometheus.GaugeOpts(prometheus.Opts{
	// 	Name: "dns_pingor_response_times_ns",
	// 	Help: "DNS response time in nano seconds",
	// }), []string{"url"})
)

func RegisterResult(res result.ResultInterface) {
	httpSuccessRate.Set(float64(res.GetSuccessRate()))

	iterateSubResults(res, func(res result.ResultInterface) {
		switch res.(type) {
		case http.Result:
			registerHTTPResult(res.(http.Result))
		}
	})
}

func registerHTTPResult(res http.Result) {
	httpResponsesTimes.
		WithLabelValues(res.URL).
		Set(float64(res.GetTime().Nanoseconds()))
}

func iterateSubResults(res result.ResultInterface, callback func(res result.ResultInterface)) {
	for _, innerRes := range res.GetSubResults() {
		callback(innerRes)
		if len(innerRes.GetSubResults()) > 0 {
			iterateSubResults(innerRes, callback)
		}
	}
}
