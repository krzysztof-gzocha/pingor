package log

import "github.com/Sirupsen/logrus"

// LogrusWrapper will wrap Logrus library to be easily testable
type LogrusWrapper struct {
	logger logrus.FieldLogger
}

// NewLogrusWrapper will return a wrapper around logrus library, so it can be easily used across the library
func NewLogrusWrapper(logger logrus.FieldLogger) *LogrusWrapper {
	return &LogrusWrapper{
		logger: logger,
	}
}

// Errorf will call underlaying logrus.Errorf
func (l *LogrusWrapper) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Debugf will call underlaying logrus.Debugf
func (l *LogrusWrapper) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Infof will call underlaying logrus.Infof
func (l *LogrusWrapper) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// WithField will call underlaying logrus.WithField
func (l *LogrusWrapper) WithField(format string, args interface{}) Logger {
	return NewLogrusWrapper(l.logger.WithField(format, args))
}
