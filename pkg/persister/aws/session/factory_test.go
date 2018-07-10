// +build unit

package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
	sess, err := CreateSession("some-region")
	assert.Nil(t, err)
	assert.NotNil(t, sess)
}
