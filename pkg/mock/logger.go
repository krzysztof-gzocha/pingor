package mock

import (
	"github.com/krzysztof-gzocha/pingor/pkg/log"
	"github.com/stretchr/testify/mock"
)

type Logger struct {
	mock.Mock
}

func (m *Logger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *Logger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *Logger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *Logger) WithField(key string, value interface{}) log.LoggerInterface {
	m.Called(key, value)

	return m
}