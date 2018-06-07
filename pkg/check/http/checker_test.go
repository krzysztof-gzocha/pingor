// +build unit

package http

import (
	"testing"

	"net/http"

	"context"

	"errors"

	"github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewChecker(t *testing.T) {
	client := mock.HttpClientMock{}
	c := NewChecker(client, "https://google.com")

	assert.NotNil(t, c)
	assert.IsType(t, Checker{}, c)
}

func TestChecker_Check_Success(t *testing.T) {
	client := mock.HttpClientMock{}
	client.
		On("Get", "google.com").
		Once().
		Return(&http.Response{StatusCode: http.StatusOK}, nil)

	client.
		On("Get", "wp.pl").
		Once().
		Return(&http.Response{StatusCode: http.StatusOK}, nil)

	c := NewChecker(client, "google.com", "wp.pl")
	result := c.Check(context.TODO())

	assert.True(t, result.IsSuccess())
	assert.Equal(t, float32(1), result.GetSuccessRate())
	assert.NotEmpty(t, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 2)
}

func TestChecker_Check_BadStatusCode(t *testing.T) {
	client := mock.HttpClientMock{}
	client.
		On("Get", "google.com").
		Once().
		Return(&http.Response{StatusCode: http.StatusOK}, nil)

	client.
		On("Get", "wp.pl").
		Once().
		Return(&http.Response{StatusCode: http.StatusNotFound}, nil)

	c := NewChecker(client, "google.com", "wp.pl")
	result := c.Check(context.TODO())

	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0.5), result.GetSuccessRate())
	assert.NotEmpty(t, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 2)
}

func TestChecker_Check_ErrorClient(t *testing.T) {
	client := mock.HttpClientMock{}
	client.
		On("Get", "google.com").
		Once().
		Return(&http.Response{StatusCode: http.StatusOK}, nil)

	client.
		On("Get", "wp.pl").
		Once().
		Return(&http.Response{StatusCode: http.StatusNotFound}, errors.New("client err"))

	c := NewChecker(client, "google.com", "wp.pl")
	result := c.Check(context.TODO())

	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0.5), result.GetSuccessRate())
	assert.NotEmpty(t, result.GetTime())
	assert.NotEmpty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 2)
}

func TestChecker_Check_NoUrlProvided(t *testing.T) {
	client := mock.HttpClientMock{}
	c := NewChecker(client)
	result := c.Check(context.TODO())

	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0), result.GetSuccessRate())
	assert.Empty(t, result.GetTime())
	assert.Empty(t, result.GetMessage())
	assert.Len(t, result.GetSubResults(), 0)
}
