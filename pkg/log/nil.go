package log

type Nil struct{}

func (n *Nil) Errorf(format string, args ...interface{}) {}
func (n *Nil) Infof(format string, args ...interface{})  {}
func (n *Nil) Debugf(format string, args ...interface{}) {}
func (n *Nil) WithField(key string, value interface{}) LoggerInterface {
	return n
}
