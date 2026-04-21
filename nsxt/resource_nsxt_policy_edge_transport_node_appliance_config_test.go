//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitNsxt_getApplianceConfigFromSchema_nilOrEmptyReturnsNil(t *testing.T) {
	require.Nil(t, getApplianceConfigFromSchema(nil))
	require.Nil(t, getApplianceConfigFromSchema([]interface{}{}))
}

func TestUnitNsxt_getApplianceConfigFromSchema_dnsServers(t *testing.T) {
	cfg := []interface{}{
		map[string]interface{}{
			"allow_ssh_root_login": false,
			"enable_ssh":           false,
			"enable_upt_mode":      false,
			"dns_servers":          []interface{}{"10.0.0.1", "10.0.0.2"},
		},
	}
	out := getApplianceConfigFromSchema(cfg)
	require.NotNil(t, out)
	require.Equal(t, []string{"10.0.0.1", "10.0.0.2"}, out.DnsServers)
}

func TestUnitNsxt_getApplianceConfigFromSchema_optionalDnsOmitted(t *testing.T) {
	cfg := []interface{}{
		map[string]interface{}{
			"allow_ssh_root_login": true,
			"enable_ssh":           true,
			"enable_upt_mode":      false,
		},
	}
	out := getApplianceConfigFromSchema(cfg)
	require.NotNil(t, out)
	require.Empty(t, out.DnsServers)
}

func TestUnitNsxt_getApplianceConfigFromSchema_syslogServerSchemaKey(t *testing.T) {
	cfg := []interface{}{
		map[string]interface{}{
			"allow_ssh_root_login": false,
			"enable_ssh":           false,
			"enable_upt_mode":      false,
			"syslog_server": []interface{}{
				map[string]interface{}{
					"log_level": "INFO",
					"port":      "514",
					"protocol":  "UDP",
					"server":    "10.1.1.1",
				},
			},
		},
	}
	out := getApplianceConfigFromSchema(cfg)
	require.NotNil(t, out)
	require.Len(t, out.SyslogServers, 1)
	require.NotNil(t, out.SyslogServers[0].Server)
	require.Equal(t, "10.1.1.1", *out.SyslogServers[0].Server)
}
