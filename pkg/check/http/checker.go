package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	pkgHttp "github.com/krzysztof-gzocha/pingor/pkg/http"
	"github.com/krzysztof-gzocha/pingor/pkg/log"
)

// Checker will make HTTP request to provided URLs and will return positive result if HTTP status will be 200 OK
type Checker struct {
	logger     log.LoggerInterface
	httpClient pkgHttp.ClientInterface
	urls       []string
}

// NewChecker will return new instance of Checker
func NewChecker(logger log.LoggerInterface, httpClient pkgHttp.ClientInterface, urls ...string) Checker {
	return Checker{logger: logger, httpClient: httpClient, urls: urls}
}

// Check will send HTTP request to all provided URLs and check HTTP statuses of the response.
// Status code have to be "200" to be recognized as success.
func (c Checker) Check(ctx context.Context) result.ResultInterface {
	if len(c.urls) == 0 {
		return result.Result{}
	}

	overallResult := result.Result{
		Success: true,
		Message: fmt.Sprintf("Checking HTTP status for %d URLs", len(c.urls)),
	}

	for _, url := range c.urls {
		overallResult.SubResults = append(overallResult.SubResults, c.singleCheck(ctx, url))
	}

	for _, res := range overallResult.SubResults {
		if !res.IsSuccess() {
			overallResult.Success = false
		}

		overallResult.Time += res.GetTime()
		overallResult.SuccessRate += res.GetSuccessRate()
	}

	overallResult.Time /= time.Duration(len(overallResult.SubResults))
	overallResult.SuccessRate /= float32(len(overallResult.SubResults))
	c.logger.
		WithField("successRate", overallResult.SuccessRate*100).
		Debugf("HTTP check is done")

	return overallResult
}

func (c Checker) singleCheck(ctx context.Context, url string) result.ResultInterface {
	c.logger.WithField("url", url).Debugf("Starting to check for HTTP status")

	res := result.Result{Success: true}
	start := time.Now()
	resp, err := c.httpClient.Get(url)
	diff := time.Now().Sub(start)
	if err != nil {
		res.Message = fmt.Sprintf("%T:%s: Failed to get URL: %s", c, url, err.Error())
		res.Success = false
		res.SuccessRate = 0

		return res
	}

	if resp.StatusCode != http.StatusOK {
		res.Success = false
		res.Message = fmt.Sprintf("%T:%s: Expecting status 200, but got %d", c, url, resp.StatusCode)
		res.SuccessRate = 0
		res.Time = diff

		return res
	}

	res.Message = fmt.Sprintf("%T:%s: Status is 200", c, url)
	res.SuccessRate = 1
	res.Time = diff

	return res
}
