//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitNsxt_resourceNsxtPolicyBgpConfigImporter(t *testing.T) {
	res := resourceNsxtPolicyBgpConfig()

	t.Run("valid policy path is parsed into gateway and locale service", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-0s/t0-1/locale-services/default")

		rd, err := resourceNsxtPolicyBgpConfigImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/tier-0s/t0-1", rd[0].Get("gateway_path"))
		assert.Equal(t, "t0-1", rd[0].Get("gateway_id"))
		assert.Equal(t, "default", rd[0].Get("locale_service_id"))
		assert.NotEmpty(t, rd[0].Id())
	})

	t.Run("path without locale-services segment fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-0s/t0-1")

		_, err := resourceNsxtPolicyBgpConfigImporter(d, nil)
		require.Error(t, err)
	})
}

func TestUnitNsxt_nsxtGatewayFloodProtectionProfileBindingImporter(t *testing.T) {
	res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()

	t.Run("valid path splits parent path and binding id", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-0s/t0-1/flood-protection-profile-bindings/binding-1")

		rd, err := nsxtGatewayFloodProtectionProfileBindingImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/tier-0s/t0-1", rd[0].Get("parent_path"))
		assert.Equal(t, "binding-1", rd[0].Id())
	})

	t.Run("path missing the bindings segment fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-0s/t0-1")

		_, err := nsxtGatewayFloodProtectionProfileBindingImporter(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid importID")
	})
}

func TestUnitNsxt_resourceNsxtPolicyIPAddressAllocationImport(t *testing.T) {
	res := resourceNsxtPolicyIPAddressAllocation()

	t.Run("policy path is parsed into pool path and allocation id", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/ip-pools/my-pool/ip-allocations/alloc-1")

		rd, err := resourceNsxtPolicyIPAddressAllocationImport(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/ip-pools/my-pool", rd[0].Get("pool_path"))
		assert.Equal(t, "alloc-1", rd[0].Id())
	})

	t.Run("legacy format with wrong segment count fails without contacting NSX", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("not-a-policy-path")

		_, err := resourceNsxtPolicyIPAddressAllocationImport(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Please provide")
	})
}

func TestUnitNsxt_resourceNsxtPolicyL2VpnServiceImport(t *testing.T) {
	res := resourceNsxtPolicyL2VpnService()

	t.Run("path with locale service sets locale_service_path", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-0s/t0-1/locale-services/default/l2vpn-services/svc-1")

		rd, err := resourceNsxtPolicyL2VpnServiceImport(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "svc-1", rd[0].Id())
		assert.Equal(t, "/infra/tier-0s/t0-1/locale-services/default", rd[0].Get("locale_service_path"))
	})

	t.Run("path without locale service sets gateway_path", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-0s/t0-1/l2vpn-services/svc-1")

		rd, err := resourceNsxtPolicyL2VpnServiceImport(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "svc-1", rd[0].Id())
		assert.Equal(t, "/infra/tier-0s/t0-1", rd[0].Get("gateway_path"))
	})

	t.Run("unexpected segment count fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-0s/t0-1")

		_, err := resourceNsxtPolicyL2VpnServiceImport(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Expected policy path")
	})
}

func TestUnitNsxt_resourceNsxtPolicyEvpnConfigImport(t *testing.T) {
	res := resourceNsxtPolicyEvpnConfig()

	t.Run("tier-0 gateway path sets gateway_path and the gateway id", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-0s/t0-1")

		rd, err := resourceNsxtPolicyEvpnConfigImport(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/tier-0s/t0-1", rd[0].Get("gateway_path"))
		assert.Equal(t, "t0-1", rd[0].Id())
	})
}

func TestUnitNsxt_resourceNsxtPolicyGatewayRedistributionConfigImport(t *testing.T) {
	res := resourceNsxtPolicyGatewayRedistributionConfig()

	t.Run("wrong segment count fails without contacting NSX", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("only-one-segment")

		_, err := resourceNsxtPolicyGatewayRedistributionConfigImport(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Please provide")
	})
}

func TestUnitNsxt_resourceNsxtPolicyIPPoolSubnetImport(t *testing.T) {
	res := resourceNsxtPolicyIPPoolBlockSubnet()

	t.Run("policy path is parsed into pool path and subnet id", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/ip-pools/my-pool/ip-subnets/subnet-1")

		rd, err := resourceNsxtPolicyIPPoolSubnetImport(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/ip-pools/my-pool", rd[0].Get("pool_path"))
		assert.Equal(t, "subnet-1", rd[0].Id())
	})

	t.Run("malformed id fails without contacting NSX", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("not-a-policy-path")

		_, err := resourceNsxtPolicyIPPoolSubnetImport(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Please provide")
	})
}

func TestUnitNsxt_nsxtDistributedFloodProtectionProfileBindingImporter(t *testing.T) {
	res := resourceNsxtPolicyDistributedFloodProtectionProfileBinding()

	t.Run("valid path splits parent group path and binding id", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/domains/default/groups/g1/firewall-flood-protection-profile-binding-maps/binding-1")

		rd, err := nsxtDistributedFloodProtectionProfileBindingImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/domains/default/groups/g1", rd[0].Get("group_path"))
	})

	t.Run("path missing the binding-maps segment fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/domains/default/groups/g1")

		_, err := nsxtDistributedFloodProtectionProfileBindingImporter(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid importID")
	})
}

func TestUnitNsxt_resourceNsxtPolicyGatewayDNSForwarderImport(t *testing.T) {
	res := resourceNsxtPolicyGatewayDNSForwarder()

	t.Run("tier-0 gateway path sets gateway_path and the gateway id", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-0s/t0-1")

		rd, err := resourceNsxtPolicyGatewayDNSForwarderImport(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/tier-0s/t0-1", rd[0].Get("gateway_path"))
		assert.Equal(t, "t0-1", rd[0].Id())
	})

	t.Run("empty id fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("")

		_, err := resourceNsxtPolicyGatewayDNSForwarderImport(d, nil)
		require.Error(t, err)
	})
}

func TestUnitNsxt_resourceNsxtPolicyEvpnTunnelEndpointImport(t *testing.T) {
	res := resourceNsxtPolicyEvpnTunnelEndpoint()

	t.Run("four-segment id is parsed into gateway/locale-service/interface/endpoint", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("t0-1/default/iface-1/endpoint-1")

		rd, err := resourceNsxtPolicyEvpnTunnelEndpointImport(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "t0-1", rd[0].Get("gateway_id"))
		assert.Equal(t, "default", rd[0].Get("locale_service_id"))
		assert.Equal(t, "/infra/tier-0s/t0-1/locale-services/default/interfaces/iface-1", rd[0].Get("external_interface_path"))
		assert.Equal(t, "endpoint-1", rd[0].Id())
	})

	t.Run("wrong segment count fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("t0-1/default")

		_, err := resourceNsxtPolicyEvpnTunnelEndpointImport(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Please provide")
	})
}

func TestUnitNsxt_nsxtGatewayResourceImporter(t *testing.T) {
	res := resourceNsxtPolicyFixedSegment()

	t.Run("policy path sets connectivity_path from the segments path", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-1s/gw1/segments/seg-1")

		rd, err := nsxtGatewayResourceImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/tier-1s/gw1", rd[0].Get("connectivity_path"))
	})

	t.Run("legacy gatewayID/segmentID format builds the connectivity path", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("gw1/seg-1")

		rd, err := nsxtGatewayResourceImporter(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "seg-1", rd[0].Id())
		assert.Equal(t, "/infra/tier-1s/gw1", rd[0].Get("connectivity_path"))
	})

	t.Run("single-segment legacy id fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("only-one-segment")

		_, err := nsxtGatewayResourceImporter(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Import format")
	})
}

func TestUnitNsxt_resourceNsxtPolicyComputeSubClusterImporter(t *testing.T) {
	res := resourceNsxtPolicyComputeSubCluster()

	t.Run("valid path extracts enforcement point and site path", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/global-infra/sites/site1/enforcement-points/ep1/sub-clusters/sc1")

		rd, err := resourceNsxtPolicyComputeSubClusterImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "ep1", rd[0].Get("enforcement_point"))
		assert.Equal(t, "/global-infra/sites/site1", rd[0].Get("site_path"))
		assert.Equal(t, "sc1", rd[0].Id())
	})

	t.Run("path without the sub-clusters segment fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/global-infra/sites/site1/enforcement-points/ep1")

		_, err := resourceNsxtPolicyComputeSubClusterImporter(d, nil)
		require.Error(t, err)
	})
}

func TestUnitNsxt_resourceNsxtPolicyEdgeHighAvailabilityProfileImporter(t *testing.T) {
	res := resourceNsxtPolicyEdgeHighAvailabilityProfile()

	t.Run("valid path extracts enforcement point and site path", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/global-infra/sites/site1/enforcement-points/ep1/edge-cluster-high-availability-profiles/prof1")

		rd, err := resourceNsxtPolicyEdgeHighAvailabilityProfileImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "ep1", rd[0].Get("enforcement_point"))
		assert.Equal(t, "/global-infra/sites/site1", rd[0].Get("site_path"))
		assert.Equal(t, "prof1", rd[0].Id())
	})
}

func TestUnitNsxt_nsxtSecurityPolicyContainerClusterImporter(t *testing.T) {
	res := resourceNsxtPolicySecurityPolicyContainerCluster()

	t.Run("valid path extracts the policy path", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/domains/default/security-policies/sp1/container-cluster-span")

		rd, err := nsxtSecurityPolicyContainerClusterImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/domains/default/security-policies/sp1", rd[0].Get("policy_path"))
	})

	t.Run("path without container-cluster-span segment fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/domains/default/security-policies/sp1")

		_, err := nsxtSecurityPolicyContainerClusterImporter(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid path")
	})
}

func TestUnitNsxt_nsxtIntrusionServiceGatewayPolicyRuleImporter(t *testing.T) {
	res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()

	t.Run("valid path extracts the policy path", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/domains/default/gateway-policies/gp1/rules/rule1")

		rd, err := nsxtIntrusionServiceGatewayPolicyRuleImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/domains/default/gateway-policies/gp1", rd[0].Get("policy_path"))
	})

	t.Run("path without a rule segment fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/domains/default/gateway-policies/gp1")

		_, err := nsxtIntrusionServiceGatewayPolicyRuleImporter(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid path")
	})
}

func TestUnitNsxt_nsxtVPNServiceResourceImporter(t *testing.T) {
	res := resourceNsxtPolicyIPSecVpnLocalEndpoint()

	t.Run("valid path sets service_path from the local-endpoints path", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-1s/aaa/locale-services/default/ipsec-vpn-services/bbb/local-endpoints/ccc")

		rd, err := nsxtVPNServiceResourceImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "ccc", rd[0].Id())
		assert.Equal(t, "/infra/tier-1s/aaa/locale-services/default/ipsec-vpn-services/bbb", rd[0].Get("service_path"))
	})

	t.Run("too few segments fails", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/tier-1s/aaa")

		_, err := nsxtVPNServiceResourceImporter(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "VPN Local Endpoint Path expected")
	})
}

func TestUnitNsxt_nsxtIntrusionServicePolicyRuleImporter(t *testing.T) {
	res := resourceNsxtPolicyIntrusionServicePolicyRule()

	t.Run("valid path extracts the policy path", func(t *testing.T) {
		d := res.TestResourceData()
		d.SetId("/infra/domains/default/security-policies/sp1/rules/rule1")

		rd, err := nsxtIntrusionServicePolicyRuleImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "/infra/domains/default/security-policies/sp1", rd[0].Get("policy_path"))
	})
}
