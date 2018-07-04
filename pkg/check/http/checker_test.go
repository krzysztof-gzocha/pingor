// +build unit

package http

import (
	"testing"

	"net/http"

	"context"

	"errors"

	pkgMock "github.com/krzysztof-gzocha/pingor/pkg/mock"
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
	logger.On("WithField", "url", "wp.pl")
	logger.On("WithField", "successRate", mock.Anything)
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

	c := NewChecker(logger, client, "google.com", "wp.pl")
	result := c.Check(context.TODO())

	assert.True(t, result.IsSuccess())
	assert.Equal(t, float32(1), result.GetSuccessRate())
	assert.NotEmpty(t, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 2)
	logger.AssertExpectations(t)
}

func TestChecker_Check_BadStatusCode(t *testing.T) {
	logger := &pkgMock.Logger{}
	logger.On("WithField", "url", "google.com")
	logger.On("WithField", "url", "wp.pl")
	logger.On("WithField", "successRate", mock.Anything)
	logger.On("Debugf", mock.Anything, mock.Anything)
	client := pkgMock.HttpClientMock{}
	client.
		On("Get", "google.com").
		Once().
		Return(&http.Response{StatusCode: http.StatusOK}, nil)

	client.
		On("Get", "wp.pl").
		Once().
		Return(&http.Response{StatusCode: http.StatusNotFound}, nil)

	c := NewChecker(logger, client, "google.com", "wp.pl")
	result := c.Check(context.TODO())

	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0.5), result.GetSuccessRate())
	assert.NotEmpty(t, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 2)
	logger.AssertExpectations(t)
}

func TestChecker_Check_ErrorClient(t *testing.T) {
	logger := &pkgMock.Logger{}
	logger.On("WithField", "url", "google.com")
	logger.On("WithField", "url", "wp.pl")
	logger.On("WithField", "successRate", mock.Anything)
	logger.On("Debugf", mock.Anything, mock.Anything)
	client := pkgMock.HttpClientMock{}
	client.
		On("Get", "google.com").
		Once().
		Return(&http.Response{StatusCode: http.StatusOK}, nil)

	client.
		On("Get", "wp.pl").
		Once().
		Return(&http.Response{StatusCode: http.StatusNotFound}, errors.New("client err"))

	c := NewChecker(logger, client, "google.com", "wp.pl")
	result := c.Check(context.TODO())

	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0.5), result.GetSuccessRate())
	assert.NotEmpty(t, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 2)
	logger.AssertExpectations(t)
}

func TestChecker_Check_NoUrlProvided(t *testing.T) {
	logger := &pkgMock.Logger{}
	client := pkgMock.HttpClientMock{}
	c := NewChecker(logger, client)
	result := c.Check(context.TODO())

	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0), result.GetSuccessRate())
	assert.Empty(t, result.GetTime())
	assert.Empty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 0)
	logger.AssertExpectations(t)
}
