// +build unit

package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDns_ResolveHost(t *testing.T) {
	d := Dns{}
	assert.Implements(t, (*Resolver)(nil), d)

	result, err := d.ResolveHost("google.com")
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Ip)
	assert.NotZero(t, result.Time)
}
