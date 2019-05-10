package http

import "net/http"

// Client should be implemented by any HTTP client capable to perform GET request
type Client interface {
	Get(url string) (*http.Response, error)
}
