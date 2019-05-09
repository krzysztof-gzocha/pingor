package log

// Nil is no-operation logger
type Nil struct{}

// Errorf will do nothing
func (n *Nil) Errorf(format string, args ...interface{}) {}

// Infof will do nothing
func (n *Nil) Infof(format string, args ...interface{}) {}

// Debugf will do nothing
func (n *Nil) Debugf(format string, args ...interface{}) {}

// WithField will do nothing
func (n *Nil) WithField(key string, value interface{}) LoggerInterface {
	return n
}
