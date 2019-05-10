package log

// Logger used across the application
type Logger interface {
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
}
