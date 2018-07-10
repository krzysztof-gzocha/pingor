// +build unit

package log

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogrusWrapper(t *testing.T) {
	l := NewLogrusWrapper(logrus.New())
	assert.NotNil(t, l)
}

func TestLogrusWrapper_Debugf(t *testing.T) {
	l := NewLogrusWrapper(logrus.New())
	assert.NotPanics(t, func() {
		l.Errorf("Test %s", "abc")
		l.Debugf("Test %s", "abc")
		l.Infof("Test %s", "abc")
	})
}

func TestLogrusWrapper_WithField(t *testing.T) {
	l := NewLogrusWrapper(logrus.New())

	assert.NotEqual(t, l, l.WithField("test", "test"))
}
