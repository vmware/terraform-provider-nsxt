//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestUnitNsxt_policyTier0GatewayInterfaceOspfSet(t *testing.T) {
	res := resourceNsxtPolicyTier0GatewayInterface()

	t.Run("no ospf config is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"type": model.Tier0Interface_TYPE_EXTERNAL,
		})
		obj := &model.Tier0Interface{}
		require.NoError(t, policyTier0GatewayInterfaceOspfSet(d, newGoMockProviderClient(), obj))
		assert.Nil(t, obj.Ospf)
	})

	t.Run("fails on global manager", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"type": model.Tier0Interface_TYPE_EXTERNAL,
			"ospf": []interface{}{
				map[string]interface{}{"enabled": true},
			},
		})
		obj := &model.Tier0Interface{}
		err := policyTier0GatewayInterfaceOspfSet(d, newGoMockGlobalProviderClient(), obj)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})

	t.Run("fails on non-external interface", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"type": model.Tier0Interface_TYPE_SERVICE,
			"ospf": []interface{}{
				map[string]interface{}{"enabled": true},
			},
		})
		obj := &model.Tier0Interface{}
		err := policyTier0GatewayInterfaceOspfSet(d, newGoMockProviderClient(), obj)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "EXTERNAL interface")
	})

	t.Run("populates the Ospf struct on an external interface", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"type": model.Tier0Interface_TYPE_EXTERNAL,
			"ospf": []interface{}{
				map[string]interface{}{
					"enabled":          true,
					"enable_bfd":       true,
					"bfd_profile_path": "/infra/bfd-profiles/bp1",
					"area_path":        "/infra/tier-0s/t0/ospf/areas/a1",
					"network_type":     model.PolicyInterfaceOspfConfig_NETWORK_TYPE_P2P,
					"hello_interval":   5,
					"dead_interval":    20,
				},
			},
		})
		obj := &model.Tier0Interface{}
		require.NoError(t, policyTier0GatewayInterfaceOspfSet(d, newGoMockProviderClient(), obj))
		require.NotNil(t, obj.Ospf)
		assert.True(t, *obj.Ospf.Enabled)
		assert.True(t, *obj.Ospf.EnableBfd)
		require.NotNil(t, obj.Ospf.BfdPath)
		assert.Equal(t, "/infra/bfd-profiles/bp1", *obj.Ospf.BfdPath)
		assert.Equal(t, "/infra/tier-0s/t0/ospf/areas/a1", *obj.Ospf.OspfArea)
		assert.Equal(t, model.PolicyInterfaceOspfConfig_NETWORK_TYPE_P2P, *obj.Ospf.NetworkType)
		assert.EqualValues(t, 5, *obj.Ospf.HelloInterval)
		assert.EqualValues(t, 20, *obj.Ospf.DeadInterval)
	})

	t.Run("empty bfd_profile_path leaves BfdPath nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"type": model.Tier0Interface_TYPE_EXTERNAL,
			"ospf": []interface{}{
				map[string]interface{}{"enabled": true},
			},
		})
		obj := &model.Tier0Interface{}
		require.NoError(t, policyTier0GatewayInterfaceOspfSet(d, newGoMockProviderClient(), obj))
		require.NotNil(t, obj.Ospf)
		assert.Nil(t, obj.Ospf.BfdPath)
	})
}
