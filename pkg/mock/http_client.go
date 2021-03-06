package mock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// HttpClientMock is just a mock
type HttpClientMock struct {
	mock.Mock
}

// Get is just a mock
func (m HttpClientMock) Get(url string) (*http.Response, error) {
	args := m.Called(url)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*http.Response), args.Error(1)
}
