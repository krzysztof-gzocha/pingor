package http

import (
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
)

// Result is main data transfer object used to store all the results from checkers
type Result struct {
	result.Result
	URL string `json:"url,omitempty"`
}
