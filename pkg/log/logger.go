package log

// LoggerInterface used across the application
type LoggerInterface interface {
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	WithField(key string, value interface{}) LoggerInterface
}
