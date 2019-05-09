// +build unit

package http

import (
	"context"
	"net/http"
	"testing"

	pkgMock "github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewChecker(t *testing.T) {
	client := pkgMock.HttpClientMock{}
	c := NewChecker(&pkgMock.Logger{}, client, "https://google.com")

	assert.NotNil(t, c)
	assert.IsType(t, Checker{}, c)
}

func TestChecker_Check_Success(t *testing.T) {
	logger := &pkgMock.Logger{}
	logger.On("WithField", "url", "google.com")
	logger.On("Debugf", mock.Anything, mock.Anything)
	client := pkgMock.HttpClientMock{}
	client.
		On("Get", "google.com").
		Once().
		Return(&http.Response{StatusCode: http.StatusOK}, nil)

	client.
		On("Get", "wp.pl").
		Once().
		Return(&http.Response{StatusCode: http.StatusOK}, nil)

	c := NewChecker(logger, client, "google.com")
	result := c.Check(context.TODO())

	assert.True(t, result.IsSuccess())
	assert.Equal(t, float32(1), result.GetSuccessRate())
	assert.NotEmpty(t, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	logger.AssertExpectations(t)
}

func TestChecker_Check_BadStatusCode(t *testing.T) {
	logger := &pkgMock.Logger{}
	logger.On("WithField", "url", "google.com")
	logger.On("Debugf", mock.Anything, mock.Anything)
	client := pkgMock.HttpClientMock{}
	client.
		On("Get", "google.com").
		Once().
		Return(&http.Response{StatusCode: http.StatusInternalServerError}, nil)

	c := NewChecker(logger, client, "google.com")
	result := c.Check(context.TODO())

	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0), result.GetSuccessRate())
	assert.NotEmpty(t, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	logger.AssertExpectations(t)
}

func TestChecker_Check_ErrorClient(t *testing.T) {
	logger := &pkgMock.Logger{}
	logger.On("WithField", "url", "google.com")
	logger.On("Debugf", mock.Anything, mock.Anything)
	client := pkgMock.HttpClientMock{}
	client.
		On("Get", "google.com").
		Once().
		Return(nil, errors.New("err"))

	c := NewChecker(logger, client, "google.com")
	result := c.Check(context.TODO())

	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0), result.GetSuccessRate())
	assert.Zero(t, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	logger.AssertExpectations(t)
}
