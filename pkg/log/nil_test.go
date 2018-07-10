// +build unit

package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNil_WithField(t *testing.T) {
	n := &Nil{}
	assert.Equal(t, n, n.WithField("t", "t"))
}

func TestNil_NotPanics(t *testing.T) {
	n := &Nil{}
	assert.NotPanics(t, func() {
		n.Errorf("t", "t")
		n.Debugf("t", "t")
		n.Infof("t", "t")
	})
}
