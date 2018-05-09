// +build unit

package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDns_ResolveHost(t *testing.T) {
	dns := Dns{}
	assert.Implements(t, (*DnsInterface)(nil), dns)

	result, err := dns.ResolveHost("google.com")
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Ip)
	assert.NotZero(t, result.Time)
}
