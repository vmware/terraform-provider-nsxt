// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitNsxt_getVMDeploymentConfigFromSchema_hostIdOmitted(t *testing.T) {
	cfg := []interface{}{
		map[string]interface{}{
			"compute_manager_id": "vc1",
			"compute_id":         "cluster1",
			"storage_id":         "datastore1",
			"compute_folder_id":  "",
			"host_id":            "",
		},
	}
	out, err := getVMDeploymentConfigFromSchema(cfg)
	require.NoError(t, err)
	require.NotNil(t, out)
	optHostID, err := out.Optional("host_id")
	require.NoError(t, err)
	require.False(t, optHostID.IsSet())
}

func TestUnitNsxt_getVMDeploymentConfigFromSchema_hostIdSet(t *testing.T) {
	cfg := []interface{}{
		map[string]interface{}{
			"compute_manager_id": "vc1",
			"compute_id":         "cluster1",
			"storage_id":         "datastore1",
			"compute_folder_id":  "",
			"host_id":            "host-1",
		},
	}
	out, err := getVMDeploymentConfigFromSchema(cfg)
	require.NoError(t, err)
	require.NotNil(t, out)
	optHostID, err := out.Optional("host_id")
	require.NoError(t, err)
	require.True(t, optHostID.IsSet())
	hostID, err := optHostID.StringValue()
	require.NoError(t, err)
	require.Equal(t, "host-1", hostID)
}
