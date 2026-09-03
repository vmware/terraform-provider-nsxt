// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitNsxt_normalizeL2VpnTransportTunnelPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "flat T0 path returned unchanged",
			input:    "/infra/tier-0s/t0/ipsec-vpn-services/svc1/sessions/sess1",
			expected: "/infra/tier-0s/t0/ipsec-vpn-services/svc1/sessions/sess1",
		},
		{
			name:     "locale-services default segment stripped from T0 path",
			input:    "/infra/tier-0s/t0/locale-services/default/ipsec-vpn-services/svc1/sessions/sess1",
			expected: "/infra/tier-0s/t0/ipsec-vpn-services/svc1/sessions/sess1",
		},
		{
			name:     "non-default locale-service ID stripped",
			input:    "/infra/tier-0s/t0/locale-services/ls-custom/ipsec-vpn-services/svc1/sessions/sess1",
			expected: "/infra/tier-0s/t0/ipsec-vpn-services/svc1/sessions/sess1",
		},
		{
			name:     "flat T1 path returned unchanged",
			input:    "/infra/tier-1s/t1/ipsec-vpn-services/svc1/sessions/sess1",
			expected: "/infra/tier-1s/t1/ipsec-vpn-services/svc1/sessions/sess1",
		},
		{
			name:     "locale-services default segment stripped from T1 path",
			input:    "/infra/tier-1s/t1/locale-services/default/ipsec-vpn-services/svc1/sessions/sess1",
			expected: "/infra/tier-1s/t1/ipsec-vpn-services/svc1/sessions/sess1",
		},
		{
			name:     "empty string returned unchanged",
			input:    "",
			expected: "",
		},
		{
			name:     "path ending at locale-services keyword without trailing slash unchanged",
			input:    "/infra/tier-0s/t0/locale-services",
			expected: "/infra/tier-0s/t0/locale-services",
		},
		{
			name:     "path with locale-service ID but no further segment unchanged",
			input:    "/infra/tier-0s/t0/locale-services/default",
			expected: "/infra/tier-0s/t0/locale-services/default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalizeL2VpnTransportTunnelPath(tt.input))
		})
	}
}

func TestUnitNsxt_suppressL2VpnTransportTunnelsDiff(t *testing.T) {
	flatPath := "/infra/tier-1s/t1/ipsec-vpn-services/svc1/sessions/sess1"
	localeServicePath := "/infra/tier-1s/t1/locale-services/default/ipsec-vpn-services/svc1/sessions/sess1"
	localeServicePathCustom := "/infra/tier-1s/t1/locale-services/custom/ipsec-vpn-services/svc1/sessions/sess1"
	differentPath := "/infra/tier-1s/t1/ipsec-vpn-services/svc2/sessions/sess2"

	tests := []struct {
		name     string
		k        string
		old      string
		new      string
		expected bool
	}{
		{
			name:     "count key is equal",
			k:        "transport_tunnels.#",
			old:      "1",
			new:      "1",
			expected: true,
		},
		{
			name:     "count key differs",
			k:        "transport_tunnels.#",
			old:      "1",
			new:      "2",
			expected: false,
		},
		{
			name:     "locale-service-scoped API path vs flat config path suppressed",
			k:        "transport_tunnels.0",
			old:      localeServicePath,
			new:      flatPath,
			expected: true,
		},
		{
			name:     "non-default locale-service path vs flat config path suppressed",
			k:        "transport_tunnels.0",
			old:      localeServicePathCustom,
			new:      flatPath,
			expected: true,
		},
		{
			name:     "flat path vs flat path is equal",
			k:        "transport_tunnels.0",
			old:      flatPath,
			new:      flatPath,
			expected: true,
		},
		{
			name:     "genuinely different paths are not suppressed",
			k:        "transport_tunnels.0",
			old:      flatPath,
			new:      differentPath,
			expected: false,
		},
		{
			name:     "empty old vs non-empty new is not suppressed",
			k:        "transport_tunnels.0",
			old:      "",
			new:      flatPath,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := suppressL2VpnTransportTunnelsDiff(tt.k, tt.old, tt.new, nil)
			assert.Equal(t, tt.expected, result)
		})
	}
}
