package pg

import (
	"testing"

	"gotest.tools/assert"
)

func TestNetworkTypeDetection(t *testing.T) {
	tts := []struct {
		address string
		typ     string
	}{
		{"127.0.0.1", tcp},
		{"192.168.0.21", tcp},
		{"dummy-pg-address-com", tcp},
		{"/cloudsql/dummy-connection-name", unix},
	}
	for _, tt := range tts {
		assert.Equal(t, getNetwork(tt.address), tt.typ)
	}
}
