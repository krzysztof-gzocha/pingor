package http

import "net/http"

// ClientInterface should be implemented by any HTTP client capable to perform GET request
type ClientInterface interface {
	Get(url string) (*http.Response, error)
}
