//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitNsxt_validateASPlainOrDot(t *testing.T) {
	cases := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"plain asn", "65001", false},
		{"dot notation", "1.10", false},
		{"too many dots", "1.2.3", true},
		{"non numeric plain", "abc", true},
		{"non numeric dot high", "1.abc", true},
		{"wrong type", 42, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, es := validateASPlainOrDot(tc.value, "key")
			if tc.wantErr {
				assert.NotEmpty(t, es, "expected error for %#v", tc.value)
			} else {
				assert.Empty(t, es)
			}
		})
	}
}

func TestUnitNsxt_validateASPath(t *testing.T) {
	cases := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"single asn", "65001", false},
		{"two asns", "65001 65002", false},
		{"dot in path", "1.10 65002", false},
		{"invalid token", "65001 bad", true},
		{"wrong type", 1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, es := validateASPath(tc.value, "key")
			if tc.wantErr {
				assert.NotEmpty(t, es)
			} else {
				assert.Empty(t, es)
			}
		})
	}
}

func TestUnitNsxt_validatePolicyBGPCommunity(t *testing.T) {
	cases := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"NO_EXPORT", "NO_EXPORT", false},
		{"NO_ADVERTISE", "NO_ADVERTISE", false},
		{"NO_EXPORT_SUBCONFED", "NO_EXPORT_SUBCONFED", false},
		{"aa nn", "1:2", false},
		{"aa bb nn", "1:2:3", false},
		{"too few parts", "1", true},
		{"too many parts", "1:2:3:4", true},
		{"non numeric", "a:b", true},
		{"wrong type", 99, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, es := validatePolicyBGPCommunity(tc.value, "key")
			if tc.wantErr {
				assert.NotEmpty(t, es)
			} else {
				assert.Empty(t, es)
			}
		})
	}
}

func TestUnitNsxt_validateLdapOrLdapsURL(t *testing.T) {
	v := validateLdapOrLdapsURL()
	cases := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"ldap", "ldap://ldap.example.com:389", false},
		{"ldaps", "ldaps://ldap.example.com:636", false},
		{"http rejected", "http://x", true},
		{"not a url", ":::not", true},
		{"wrong type", 1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, es := v(tc.value, "key")
			if tc.wantErr {
				assert.NotEmpty(t, es)
			} else {
				assert.Empty(t, es)
			}
		})
	}
}

func TestUnitNsxt_validateNsxtProviderHostFormat(t *testing.T) {
	v := validateNsxtProviderHostFormat()
	cases := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"host only", "nsx.example.com", false},
		{"https url", "https://nsx.example.com", false},
		{"empty invalid", "", true},
		{"wrong type", 1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, es := v(tc.value, "key")
			if tc.wantErr {
				assert.NotEmpty(t, es)
			} else {
				assert.Empty(t, es)
			}
		})
	}
}

func TestUnitNsxt_validatePortRangeAndSinglePort(t *testing.T) {
	pr := validatePortRange()
	sp := validateSinglePort()
	cases := []struct {
		name    string
		v       func(interface{}, string) ([]string, []error)
		value   interface{}
		wantErr bool
	}{
		{"range ok", pr, "80-443", false},
		{"single port ok", pr, "8080", false},
		{"range invalid", pr, "99999-100000", true},
		{"garbage", pr, "a-b", true},
		{"single ok", sp, "443", false},
		{"single bad", sp, "70000", true},
		{"single wrong type", sp, struct{}{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, es := tc.v(tc.value, "key")
			if tc.wantErr {
				assert.NotEmpty(t, es)
			} else {
				assert.Empty(t, es)
			}
		})
	}
}

func TestUnitNsxt_validatePowerOf2(t *testing.T) {
	v := validatePowerOf2(false, 64)
	cases := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"power of 2", 8, false},
		{"not pow2", 6, true},
		{"wrong type", "8", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, es := v(tc.value, "key")
			if tc.wantErr {
				assert.NotEmpty(t, es)
			} else {
				assert.Empty(t, es)
			}
		})
	}

	t.Run("allow zero", func(t *testing.T) {
		v0 := validatePowerOf2(true, 0)
		_, es := v0(0, "k")
		assert.Empty(t, es)
	})
}

func TestUnitNsxt_validateASNPairAndIPorASNPair(t *testing.T) {
	t.Run("ASN pair", func(t *testing.T) {
		_, es := validateASNPair("65001:100", "k")
		assert.Empty(t, es)
		_, es = validateASNPair("bad", "k")
		assert.NotEmpty(t, es)
		_, es = validateASNPair(1, "k")
		assert.NotEmpty(t, es)
	})
	t.Run("IP or ASN pair", func(t *testing.T) {
		_, es := validateIPorASNPair("192.0.2.1:42", "k")
		assert.Empty(t, es)
		_, es = validateIPorASNPair("65001:100", "k")
		assert.Empty(t, es)
		_, es = validateIPorASNPair("nope", "k")
		assert.NotEmpty(t, es)
	})
}

func TestUnitNsxt_validateVLANIdAndRange(t *testing.T) {
	t.Run("vlan id int", func(t *testing.T) {
		_, es := validateVLANId(100, "k")
		assert.Empty(t, es)
		_, es = validateVLANId(5000, "k")
		assert.NotEmpty(t, es)
	})
	t.Run("vlan id string", func(t *testing.T) {
		_, es := validateVLANId("4095", "k")
		assert.Empty(t, es)
	})
	t.Run("vlan range", func(t *testing.T) {
		_, es := validateVLANIdOrRange("10-20", "k")
		assert.Empty(t, es)
		_, es = validateVLANIdOrRange("1-2-3", "k")
		assert.NotEmpty(t, es)
	})
}

func TestUnitNsxt_validateStringIntBetween(t *testing.T) {
	v := validateStringIntBetween(1, 5)
	cases := []struct {
		val     string
		wantErr bool
	}{
		{"3", false},
		{"0", true},
		{"6", true},
		{"x", true},
	}
	for _, tc := range cases {
		t.Run(tc.val, func(t *testing.T) {
			_, es := v(tc.val, "k")
			if tc.wantErr {
				assert.NotEmpty(t, es)
			} else {
				assert.Empty(t, es)
			}
		})
	}
	_, es := v(1, "k")
	require.NotEmpty(t, es)
}

func TestUnitNsxt_validateSingleIPOrHostName(t *testing.T) {
	v := validateSingleIPOrHostName()
	cases := []struct {
		val     string
		wantErr bool
	}{
		{"192.0.2.1", false},
		{"host.example.com", false},
		{"bad host!", true},
	}
	for _, tc := range cases {
		t.Run(tc.val, func(t *testing.T) {
			_, es := v(tc.val, "k")
			if tc.wantErr {
				assert.NotEmpty(t, es)
			} else {
				assert.Empty(t, es)
			}
		})
	}
}

func TestUnitNsxt_validatePolicyPathAndIDAndProjectID(t *testing.T) {
	vp := validatePolicyPath()
	_, es := vp("/infra/domains/default/groups/g1", "k")
	assert.Empty(t, es)
	_, es = vp("not-a-path", "k")
	assert.NotEmpty(t, es)

	vid := validateID()
	_, es = vid("ok-id", "k")
	assert.Empty(t, es)
	_, es = vid("bad/id", "k")
	assert.NotEmpty(t, es)

	vpdef := validateProjectID(true)
	_, es = vpdef("default", "k")
	assert.Empty(t, es)

	vpno := validateProjectID(false)
	_, es = vpno("default", "k")
	assert.NotEmpty(t, es)
}

func TestUnitNsxt_validatePolicySourceDestinationGroups(t *testing.T) {
	v := validatePolicySourceDestinationGroups()
	cases := []struct {
		val     string
		wantErr bool
	}{
		{"/infra/domains/default/groups/g1", false},
		{"192.0.2.0/24", false},
		{"192.0.2.1-192.0.2.5", false},
		{"not-valid", true},
	}
	for _, tc := range cases {
		t.Run(strings.Trim(tc.val, "/"), func(t *testing.T) {
			_, es := v(tc.val, "k")
			if tc.wantErr {
				assert.NotEmpty(t, es)
			} else {
				assert.Empty(t, es)
			}
		})
	}
}

func TestUnitNsxt_validateCidrIPCidr(t *testing.T) {
	vc := validateCidr()
	_, es := vc("10.0.0.0/8", "k")
	assert.Empty(t, es)
	_, es = vc("not-cidr", "k")
	assert.NotEmpty(t, es)

	vipc := validateIPCidr()
	_, es = vipc("10.0.0.1/32", "k")
	assert.Empty(t, es)
}

func TestUnitNsxt_validateLocaleServiceAndGatewayPolicyPath(t *testing.T) {
	vl := validateLocaleServicePolicyPath()
	_, es := vl("/infra/tier-0s/mygw/locale-services/default", "k")
	assert.Empty(t, es)
	_, es = vl("/not/locale/services", "k")
	assert.NotEmpty(t, es)
	_, es = vl(42, "k")
	assert.NotEmpty(t, es)

	vg := validateGatewayPolicyPath()
	_, es = vg("/infra/tier-0s/gw1", "k")
	assert.Empty(t, es)
	_, es = vg("/infra/tier-1s/t1", "k")
	assert.Empty(t, es)
	_, es = vg("/bad", "k")
	assert.NotEmpty(t, es)
	_, es = vg(1, "k")
	assert.NotEmpty(t, es)
}

func TestUnitNsxt_isT0Gw(t *testing.T) {
	ok, err := isT0Gw("/infra/tier-0s/t0")
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = isT0Gw("/infra/tier-1s/t1")
	require.NoError(t, err)
	assert.False(t, ok)

	_, err = isT0Gw("/short")
	require.Error(t, err)
}
