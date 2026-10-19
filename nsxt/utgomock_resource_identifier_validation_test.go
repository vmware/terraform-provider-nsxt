//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

// A caller-supplied NSX identifier must be validated before any backend mutation.
func TestUnitNsxt_CreateRejectsInvalidNSXIdentifiers(t *testing.T) {
	oldVersion := util.NsxVersion
	t.Cleanup(func() { util.NsxVersion = oldVersion })
	util.NsxVersion = "9.2.0"
	for _, tc := range []struct {
		name     string
		resource func() *schema.Resource
		create   func(*schema.ResourceData, interface{}) error
	}{
		{"resourceNsxtPolicyBgpNeighborCreate", resourceNsxtPolicyBgpNeighbor, resourceNsxtPolicyBgpNeighborCreate},
		{"resourceNsxtPolicyComputeSubClusterCreate", resourceNsxtPolicyComputeSubCluster, resourceNsxtPolicyComputeSubClusterCreate},
		{"resourceNsxtPolicyConstraintCreate", resourceNsxtPolicyConstraint, resourceNsxtPolicyConstraintCreate},
		{"resourceNsxtPolicyContextProfileCreate", resourceNsxtPolicyContextProfile, resourceNsxtPolicyContextProfileCreate},
		{"resourceNsxtPolicyDhcpRelayConfigCreate", resourceNsxtPolicyDhcpRelayConfig, resourceNsxtPolicyDhcpRelayConfigCreate},
		{"resourceNsxtPolicyDhcpServerCreate", resourceNsxtPolicyDhcpServer, resourceNsxtPolicyDhcpServerCreate},
		{"resourceNsxtPolicyDhcpV4StaticBindingCreate", resourceNsxtPolicyDhcpV4StaticBinding, resourceNsxtPolicyDhcpV4StaticBindingCreate},
		{"resourceNsxtPolicyDhcpV6StaticBindingCreate", resourceNsxtPolicyDhcpV6StaticBinding, resourceNsxtPolicyDhcpV6StaticBindingCreate},
		{"resourceNsxtPolicyDistributedFloodProtectionProfileCreate", resourceNsxtPolicyDistributedFloodProtectionProfile, resourceNsxtPolicyDistributedFloodProtectionProfileCreate},
		{"resourceNsxtPolicyDistributedFloodProtectionProfileBindingCreate", resourceNsxtPolicyDistributedFloodProtectionProfileBinding, resourceNsxtPolicyDistributedFloodProtectionProfileBindingCreate},
		{"resourceNsxtPolicyDistributedVlanConnectionCreate", resourceNsxtPolicyDistributedVlanConnection, resourceNsxtPolicyDistributedVlanConnectionCreate},
		{"resourceNsxtPolicyDistributedVxlanConnectionCreate", resourceNsxtPolicyDistributedVxlanConnection, resourceNsxtPolicyDistributedVxlanConnectionCreate},
		{"resourceNsxtPolicyDNSForwarderZoneCreate", resourceNsxtPolicyDNSForwarderZone, resourceNsxtPolicyDNSForwarderZoneCreate},
		{"resourceNsxtPolicyDnsRecordCreate", resourceNsxtPolicyDnsRecord, resourceNsxtPolicyDnsRecordCreate},
		{"resourceNsxtPolicyDnsRecordAutoConfigCreate", resourceNsxtPolicyDnsRecordAutoConfig, resourceNsxtPolicyDnsRecordAutoConfigCreate},
		{"resourceNsxtPolicyDnsServiceCreate", resourceNsxtPolicyDnsService, resourceNsxtPolicyDnsServiceCreate},
		{"resourceNsxtPolicyDomainCreate", resourceNsxtPolicyDomain, resourceNsxtPolicyDomainCreate},
		{"resourceNsxtPolicyEdgeClusterCreate", resourceNsxtPolicyEdgeCluster, resourceNsxtPolicyEdgeClusterCreate},
		{"resourceNsxtPolicyEdgeHighAvailabilityProfileCreate", resourceNsxtPolicyEdgeHighAvailabilityProfile, resourceNsxtPolicyEdgeHighAvailabilityProfileCreate},
		{"resourceNsxtPolicyEdgeTransportNodeCreate", resourceNsxtPolicyEdgeTransportNode, resourceNsxtPolicyEdgeTransportNodeCreate},
		{"resourceNsxtPolicyEvpnTenantCreate", resourceNsxtPolicyEvpnTenant, resourceNsxtPolicyEvpnTenantCreate},
		{"resourceNsxtPolicyEvpnTunnelEndpointCreate", resourceNsxtPolicyEvpnTunnelEndpoint, resourceNsxtPolicyEvpnTunnelEndpointCreate},
		{"resourceNsxtPolicyGatewayCommunityListCreate", resourceNsxtPolicyGatewayCommunityList, resourceNsxtPolicyGatewayCommunityListCreate},
		{"resourceNsxtPolicyGatewayConnectionCreate", resourceNsxtPolicyGatewayConnection, resourceNsxtPolicyGatewayConnectionCreate},
		{"resourceNsxtPolicyGatewayFloodProtectionProfileCreate", resourceNsxtPolicyGatewayFloodProtectionProfile, resourceNsxtPolicyGatewayFloodProtectionProfileCreate},
		{"resourceNsxtPolicyGatewayPolicyRuleCreate", resourceNsxtPolicyGatewayRulePolicy, resourceNsxtPolicyGatewayPolicyRuleCreate},
		{"resourceNsxtPolicyGatewayPrefixListCreate", resourceNsxtPolicyGatewayPrefixList, resourceNsxtPolicyGatewayPrefixListCreate},
		{"resourceNsxtPolicyGatewayQosProfileCreate", resourceNsxtPolicyGatewayQosProfile, resourceNsxtPolicyGatewayQosProfileCreate},
		{"resourceNsxtPolicyGatewayRouteMapCreate", resourceNsxtPolicyGatewayRouteMap, resourceNsxtPolicyGatewayRouteMapCreate},
		{"resourceNsxtPolicyGlobalManagerCreate", resourceNsxtPolicyGlobalManager, resourceNsxtPolicyGlobalManagerCreate},
		{"resourceNsxtPolicyHostTransportNodeCreate", resourceNsxtPolicyHostTransportNode, resourceNsxtPolicyHostTransportNodeCreate},
		{"resourceNsxtPolicyHostTransportNodeCollectionCreate", resourceNsxtPolicyHostTransportNodeCollection, resourceNsxtPolicyHostTransportNodeCollectionCreate},
		{"resourceNsxtPolicyHostTransportNodeProfileCreate", resourceNsxtPolicyHostTransportNodeProfile, resourceNsxtPolicyHostTransportNodeProfileCreate},
		{"resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate", resourceNsxtPolicyIntrusionServiceGatewayPolicyRule, resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate},
		{"resourceNsxtPolicyIntrusionServicePolicyRuleCreate", resourceNsxtPolicyIntrusionServicePolicyRule, resourceNsxtPolicyIntrusionServicePolicyRuleCreate},
		{"resourceNsxtPolicyIntrusionServiceProfileCreate", resourceNsxtPolicyIntrusionServiceProfile, resourceNsxtPolicyIntrusionServiceProfileCreate},
		{"resourceNsxtPolicyIPAddressAllocationCreate", resourceNsxtPolicyIPAddressAllocation, resourceNsxtPolicyIPAddressAllocationCreate},
		{"resourceNsxtPolicyIPBlockCreate", resourceNsxtPolicyIPBlock, resourceNsxtPolicyIPBlockCreate},
		{"resourceNsxtPolicyIpBlockQuotaCreate", resourceNsxtPolicyIpBlockQuota, resourceNsxtPolicyIpBlockQuotaCreate},
		{"resourceNsxtPolicyIPDiscoveryProfileCreate", resourceNsxtPolicyIPDiscoveryProfile, resourceNsxtPolicyIPDiscoveryProfileCreate},
		{"resourceNsxtPolicyIPPoolCreate", resourceNsxtPolicyIPPool, resourceNsxtPolicyIPPoolCreate},
		{"resourceNsxtPolicyIPPoolBlockSubnetCreate", resourceNsxtPolicyIPPoolBlockSubnet, resourceNsxtPolicyIPPoolBlockSubnetCreate},
		{"resourceNsxtPolicyIPPoolStaticSubnetCreate", resourceNsxtPolicyIPPoolStaticSubnet, resourceNsxtPolicyIPPoolStaticSubnetCreate},
		{"resourceNsxtPolicyIPSecVpnDpdProfileCreate", resourceNsxtPolicyIPSecVpnDpdProfile, resourceNsxtPolicyIPSecVpnDpdProfileCreate},
		{"resourceNsxtPolicyIPSecVpnIkeProfileCreate", resourceNsxtPolicyIPSecVpnIkeProfile, resourceNsxtPolicyIPSecVpnIkeProfileCreate},
		{"resourceNsxtPolicyIPSecVpnLocalEndpointCreate", resourceNsxtPolicyIPSecVpnLocalEndpoint, resourceNsxtPolicyIPSecVpnLocalEndpointCreate},
		{"resourceNsxtPolicyIPSecVpnServiceCreate", resourceNsxtPolicyIPSecVpnService, resourceNsxtPolicyIPSecVpnServiceCreate},
		{"resourceNsxtPolicyIPSecVpnSessionCreate", resourceNsxtPolicyIPSecVpnSession, resourceNsxtPolicyIPSecVpnSessionCreate},
		{"resourceNsxtPolicyIPSecVpnTunnelProfileCreate", resourceNsxtPolicyIPSecVpnTunnelProfile, resourceNsxtPolicyIPSecVpnTunnelProfileCreate},
		{"resourceNsxtPolicyL2VpnServiceCreate", resourceNsxtPolicyL2VpnService, resourceNsxtPolicyL2VpnServiceCreate},
		{"resourceNsxtPolicyL2VPNSessionCreate", resourceNsxtPolicyL2VPNSession, resourceNsxtPolicyL2VPNSessionCreate},
		{"resourceNsxtPolicyL7AccessProfileCreate", resourceNsxtPolicyL7AccessProfile, resourceNsxtPolicyL7AccessProfileCreate},
		{"resourceNsxtPolicyLBClientSslProfileCreate", resourceNsxtPolicyLBClientSslProfile, resourceNsxtPolicyLBClientSslProfileCreate},
		{"resourceNsxtPolicyLBCookiePersistenceProfileCreate", resourceNsxtPolicyLBCookiePersistenceProfile, resourceNsxtPolicyLBCookiePersistenceProfileCreate},
		{"resourceNsxtPolicyLBTcpApplicationProfileCreate", resourceNsxtPolicyLBFastTcpApplicationProfile, resourceNsxtPolicyLBTcpApplicationProfileCreate},
		{"resourceNsxtPolicyLBUdpApplicationProfileCreate", resourceNsxtPolicyLBFastUdpApplicationProfile, resourceNsxtPolicyLBUdpApplicationProfileCreate},
		{"resourceNsxtPolicyLBGenericPersistenceProfileCreate", resourceNsxtPolicyLBGenericPersistenceProfile, resourceNsxtPolicyLBGenericPersistenceProfileCreate},
		{"resourceNsxtPolicyLBHttpApplicationProfileCreate", resourceNsxtPolicyLBHttpApplicationProfile, resourceNsxtPolicyLBHttpApplicationProfileCreate},
		{"resourceNsxtPolicyLBHttpMonitorProfileCreate", resourceNsxtPolicyLBHttpMonitorProfile, resourceNsxtPolicyLBHttpMonitorProfileCreate},
		{"resourceNsxtPolicyLBHttpsMonitorProfileCreate", resourceNsxtPolicyLBHttpsMonitorProfile, resourceNsxtPolicyLBHttpsMonitorProfileCreate},
		{"resourceNsxtPolicyLBIcmpMonitorProfileCreate", resourceNsxtPolicyLBIcmpMonitorProfile, resourceNsxtPolicyLBIcmpMonitorProfileCreate},
		{"resourceNsxtPolicyLBPassiveMonitorProfileCreate", resourceNsxtPolicyLBPassiveMonitorProfile, resourceNsxtPolicyLBPassiveMonitorProfileCreate},
		{"resourceNsxtPolicyLBPoolCreate", resourceNsxtPolicyLBPool, resourceNsxtPolicyLBPoolCreate},
		{"resourceNsxtPolicyLBServerSslProfileCreate", resourceNsxtPolicyLBServerSslProfile, resourceNsxtPolicyLBServerSslProfileCreate},
		{"resourceNsxtPolicyLBServiceCreate", resourceNsxtPolicyLBService, resourceNsxtPolicyLBServiceCreate},
		{"resourceNsxtPolicyLBSourceIpPersistenceProfileCreate", resourceNsxtPolicyLBSourceIpPersistenceProfile, resourceNsxtPolicyLBSourceIpPersistenceProfileCreate},
		{"resourceNsxtPolicyLBTcpMonitorProfileCreate", resourceNsxtPolicyLBTcpMonitorProfile, resourceNsxtPolicyLBTcpMonitorProfileCreate},
		{"resourceNsxtPolicyLBUdpMonitorProfileCreate", resourceNsxtPolicyLBUdpMonitorProfile, resourceNsxtPolicyLBUdpMonitorProfileCreate},
		{"resourceNsxtPolicyLBVirtualServerCreate", resourceNsxtPolicyLBVirtualServer, resourceNsxtPolicyLBVirtualServerCreate},
		{"resourceNsxtPolicyLdapIdentitySourceCreate", resourceNsxtPolicyLdapIdentitySource, resourceNsxtPolicyLdapIdentitySourceCreate},
		{"resourceNsxtPolicyMacDiscoveryProfileCreate", resourceNsxtPolicyMacDiscoveryProfile, resourceNsxtPolicyMacDiscoveryProfileCreate},
		{"resourceNsxtPolicyMetadataProxyCreate", resourceNsxtPolicyMetadataProxy, resourceNsxtPolicyMetadataProxyCreate},
		{"resourceNsxtPolicyNATRuleCreate", resourceNsxtPolicyNATRule, resourceNsxtPolicyNATRuleCreate},
		{"resourceNsxtPolicyNetworkSpanCreate", resourceNsxtPolicyNetworkSpan, resourceNsxtPolicyNetworkSpanCreate},
		{"resourceNsxtPolicyOspfAreaCreate", resourceNsxtPolicyOspfArea, resourceNsxtPolicyOspfAreaCreate},
		{"resourceNsxtPolicyProjectCreate", resourceNsxtPolicyProject, resourceNsxtPolicyProjectCreate},
		{"resourceNsxtPolicyProjectCNACreate", resourceNsxtPolicyProjectCentralizedNetworkAttachment, resourceNsxtPolicyProjectCNACreate},
		{"resourceNsxtPolicyProjectIpAddressAllocationCreate", resourceNsxtPolicyProjectIpAddressAllocation, resourceNsxtPolicyProjectIpAddressAllocationCreate},
		{"resourceNsxtPolicyQosProfileCreate", resourceNsxtPolicyQosProfile, resourceNsxtPolicyQosProfileCreate},
		{"resourceNsxtPolicyRouteControllerCreate", resourceNsxtPolicyRouteController, resourceNsxtPolicyRouteControllerCreate},
		{"resourceNsxtPolicySecurityPolicyContainerClusterCreate", resourceNsxtPolicySecurityPolicyContainerCluster, resourceNsxtPolicySecurityPolicyContainerClusterCreate},
		{"resourceNsxtPolicySecurityPolicyRuleCreate", resourceNsxtPolicySecurityPolicyRule, resourceNsxtPolicySecurityPolicyRuleCreate},
		{"resourceNsxtPolicySegmentPortCreate", resourceNsxtPolicySegmentPort, resourceNsxtPolicySegmentPortCreate},
		{"resourceNsxtPolicySegmentSecurityProfileCreate", resourceNsxtPolicySegmentSecurityProfile, resourceNsxtPolicySegmentSecurityProfileCreate},
		{"resourceNsxtPolicyServiceCreate", resourceNsxtPolicyService, resourceNsxtPolicyServiceCreate},
		{"resourceNsxtPolicyShareCreate", resourceNsxtPolicyShare, resourceNsxtPolicyShareCreate},
		{"resourceNsxtPolicySiteCreate", resourceNsxtPolicySite, resourceNsxtPolicySiteCreate},
		{"resourceNsxtPolicyStaticRouteCreate", resourceNsxtPolicyStaticRoute, resourceNsxtPolicyStaticRouteCreate},
		{"resourceNsxtPolicyStaticRouteBfdPeerCreate", resourceNsxtPolicyStaticRouteBfdPeer, resourceNsxtPolicyStaticRouteBfdPeerCreate},
		{"resourceNsxtPolicyTier0GatewayCreate", resourceNsxtPolicyTier0Gateway, resourceNsxtPolicyTier0GatewayCreate},
		{"resourceNsxtPolicyTier0GatewayInterfaceCreate", resourceNsxtPolicyTier0GatewayInterface, resourceNsxtPolicyTier0GatewayInterfaceCreate},
		{"resourceNsxtPolicyTier0InterVRFRoutingCreate", resourceNsxtPolicyTier0InterVRFRouting, resourceNsxtPolicyTier0InterVRFRoutingCreate},
		{"resourceNsxtPolicyTier1GatewayCreate", resourceNsxtPolicyTier1Gateway, resourceNsxtPolicyTier1GatewayCreate},
		{"resourceNsxtPolicyTier1GatewayInterfaceCreate", resourceNsxtPolicyTier1GatewayInterface, resourceNsxtPolicyTier1GatewayInterfaceCreate},
		{"resourceNsxtPolicyTransitGatewayCreate", resourceNsxtPolicyTransitGateway, resourceNsxtPolicyTransitGatewayCreate},
		{"resourceNsxtPolicyTransportZoneCreate", resourceNsxtPolicyTransportZone, resourceNsxtPolicyTransportZoneCreate},
		{"resourceNsxtUplinkHostSwitchProfileCreate", resourceNsxtUplinkHostSwitchProfile, resourceNsxtUplinkHostSwitchProfileCreate},
		{"resourceNsxtPolicyVniPoolCreate", resourceNsxtPolicyVniPool, resourceNsxtPolicyVniPoolCreate},
		{"resourceNsxtVtepHAHostSwitchProfileCreate", resourceNsxtVtepHAHostSwitchProfile, resourceNsxtVtepHAHostSwitchProfileCreate},
		{"resourceNsxtVpcCreate", resourceNsxtVpc, resourceNsxtVpcCreate},
		{"resourceNsxtVpcConnectivityProfileCreate", resourceNsxtVpcConnectivityProfile, resourceNsxtVpcConnectivityProfileCreate},
		{"resourceNsxtVpcEndpointCreate", resourceNsxtVpcEndpoint, resourceNsxtVpcEndpointCreate},
		{"resourceNsxtVPCGatewayPolicyCreate", resourceNsxtVPCGatewayPolicy, resourceNsxtVPCGatewayPolicyCreate},
		{"resourceNsxtVpcIpAddressAllocationCreate", resourceNsxtVpcIpAddressAllocation, resourceNsxtVpcIpAddressAllocationCreate},
		{"resourceNsxtVpcServiceEndpointCreate", resourceNsxtVpcServiceEndpoint, resourceNsxtVpcServiceEndpointCreate},
		{"resourceNsxtVpcServiceProfileCreate", resourceNsxtVpcServiceProfile, resourceNsxtVpcServiceProfileCreate},
		{"resourceNsxtVpcStaticRoutesCreate", resourceNsxtVpcStaticRoutes, resourceNsxtVpcStaticRoutesCreate},
		{"resourceNsxtVpcSubnetCreate", resourceNsxtVpcSubnet, resourceNsxtVpcSubnetCreate},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.resource()
			raw := map[string]interface{}{}
			for k, v := range map[string]interface{}{
				"nsx_id": "invalid/id", "gateway_path": "/infra/tier-0s/gateway", "segment_path": "/infra/segments/segment", "pool_path": "/infra/ip-pools/pool", "site_path": "/infra/sites/default",
				"parent_path": "/orgs/default/projects/project/vpcs/vpc", "domain": "default", "domain_id": "default", "group_path": "/infra/domains/default/groups/group",
				"locale_service_id": "default", "gateway_id": "gateway", "service_path": "/infra/tier-0s/gateway/ipsec-vpn-services/vpn", "enforcement_point": "default",
				"context": []interface{}{map[string]interface{}{"project_id": "project", "vpc_id": "vpc"}},
			} {
				if _, ok := res.Schema[k]; ok {
					raw[k] = v
				}
			}
			if !strings.Contains(tc.name, "Vpc") && !strings.Contains(tc.name, "VPC") && !strings.Contains(tc.name, "TransitGateway") && !strings.Contains(tc.name, "ProjectIPAddress") {
				delete(raw, "context")
			}
			if strings.Contains(tc.name, "Tier1") {
				raw["gateway_path"] = "/infra/tier-1s/gateway"
			}
			if tc.name == "resourceNsxtPolicyL2VPNSessionCreate" {
				raw["service_path"] = "/infra/tier-0s/gateway/l2vpn-services/vpn"
			}
			d := schema.TestResourceDataRaw(t, res.Schema, raw)
			m := newIsolatedMockProviderClient(t)
			m.Host = "https://nsx.invalid"
			if tc.name == "resourceNsxtPolicyGlobalManagerCreate" || tc.name == "resourceNsxtPolicySiteCreate" {
				m.PolicyGlobalManager = true
			}
			m.PolicyHTTPClient = &http.Client{Transport: unitPolicyTransport(func(r *http.Request) (*http.Response, error) {
				t.Fatalf("invalid identifier must be rejected before API request: %s %s", r.Method, r.URL)
				return nil, nil
			})}
			err := tc.create(d, m)
			require.ErrorContains(t, err, "nsx_id")
			require.ErrorContains(t, err, "invalid")
			require.Empty(t, d.Id())
		})
	}
}
