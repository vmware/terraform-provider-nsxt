//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
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

func TestUnitNsxt_getVMDeploymentConfigFromSchema_nil(t *testing.T) {
	out, err := getVMDeploymentConfigFromSchema(nil)
	require.NoError(t, err)
	require.Nil(t, out)
}

func TestUnitNsxt_getVMDeploymentConfigFromSchema_allFields(t *testing.T) {
	cfg := []interface{}{
		map[string]interface{}{
			"compute_manager_id": "vc1",
			"compute_id":         "cluster1",
			"storage_id":         "datastore1",
			"compute_folder_id":  "folder1",
			"host_id":            "",
			"edge_host_affinity_config": []interface{}{
				map[string]interface{}{"host_group_name": "hg1"},
			},
			"reservation_info": []interface{}{
				map[string]interface{}{
					"cpu_reservation_in_mhz":        1000,
					"cpu_reservation_in_shares":     model.CPUReservation_RESERVATION_IN_SHARES_HIGH_PRIORITY,
					"memory_reservation_percentage": 50,
				},
			},
		},
	}
	out, err := getVMDeploymentConfigFromSchema(cfg)
	require.NoError(t, err)
	require.NotNil(t, out)

	golang, errs := bindings.NewTypeConverter().ConvertToGolang(out, model.PolicyVsphereDeploymentConfigBindingType())
	require.Empty(t, errs)
	obj := golang.(model.PolicyVsphereDeploymentConfig)
	assert.Equal(t, "folder1", *obj.ComputeFolderId)
	require.NotNil(t, obj.EdgeHostAffinityConfig)
	assert.Equal(t, "hg1", *obj.EdgeHostAffinityConfig.HostGroupName)
	require.NotNil(t, obj.ReservationInfo)
	assert.Equal(t, int64(1000), *obj.ReservationInfo.CpuReservation.ReservationInMhz)
	assert.Equal(t, model.CPUReservation_RESERVATION_IN_SHARES_HIGH_PRIORITY, *obj.ReservationInfo.CpuReservation.ReservationInShares)
	assert.Equal(t, int64(50), *obj.ReservationInfo.MemoryReservation.ReservationPercentage)
}
