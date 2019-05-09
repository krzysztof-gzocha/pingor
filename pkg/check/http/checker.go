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
	url        string
}

// NewChecker will return new instance of Checker
func NewChecker(logger log.LoggerInterface, httpClient pkgHttp.ClientInterface, url string) Checker {
	return Checker{logger: logger, httpClient: httpClient, url: url}
}

// Check will send HTTP request to all provided URLs and check HTTP statuses of the response.
// Status code have to be "200" to be recognized as success.
func (c Checker) Check(ctx context.Context) result.ResultInterface {
	c.logger.WithField("url", c.url).Debugf("Starting to check for HTTP status")

	res := result.Result{Success: true, URL: c.url}
	start := time.Now()
	resp, err := c.httpClient.Get(c.url)
	diff := time.Now().Sub(start)
	if err != nil {
		res.Message = fmt.Sprintf("%T:%s: Failed to get URL: %s", c, c.url, err.Error())
		res.Success = false
		res.SuccessRate = 0

		return res
	}

	if resp.StatusCode != http.StatusOK {
		res.Success = false
		res.Message = fmt.Sprintf("%T:%s: Expecting status 200, but got %d", c, c.url, resp.StatusCode)
		res.SuccessRate = 0
		res.Time = diff

		return res
	}

	res.Message = fmt.Sprintf("%T:%s: Status is 200", c, c.url)
	res.SuccessRate = 1
	res.Time = diff

	return res
}
