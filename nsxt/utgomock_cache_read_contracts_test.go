//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

// These contracts apply to every cache-aware resource: a failed refresh must preserve
// the Terraform ID, while a confirmed backend deletion must remove it from state.
// Use the real SDK response decoding with an in-memory transport; no NSX endpoint is needed.
func TestUnitNsxt_CacheMissReadErrorContracts(t *testing.T) {
	oldCache, oldVersion := gcache, util.NsxVersion
	gcache = &typeScopedCache{byTyp: make(map[string]*resourceTypeCache)}
	t.Cleanup(func() { gcache = oldCache; util.NsxVersion = oldVersion })
	util.NsxVersion = "9.2.0"
	for _, resource := range []struct {
		name         string
		factory      func() *schema.Resource
		resourceType string
	}{
		{"policy_connectivity_policy", resourceNsxtPolicyConnectivityPolicy, resourceTypeConnectivityPolicy},
		{"policy_dhcp_v4_static_binding", resourceNsxtPolicyDhcpV4StaticBinding, resourceTypeDhcpV4StaticBindingConfig},
		{"policy_dhcp_v6_static_binding", resourceNsxtPolicyDhcpV6StaticBinding, resourceTypeDhcpV6StaticBindingConfig},
		{"policy_distributed_flood_protection_profile_binding", resourceNsxtPolicyDistributedFloodProtectionProfileBinding, resourceTypePolicyFirewallFloodProtectionProfileBindingMap},
		{"policy_gateway_policy", resourceNsxtPolicyGatewayPolicy, resourceTypeGatewayPolicy},
		{"policy_ip_address_allocation", resourceNsxtPolicyIPAddressAllocation, resourceTypeIpAddressAllocation},
		{"policy_ip_pool", resourceNsxtPolicyIPPool, resourceTypeIpAddressPool},
		{"policy_ip_pool_block_subnet", resourceNsxtPolicyIPPoolBlockSubnet, resourceTypeIpAddressPoolBlockSubnet},
		{"policy_ip_pool_static_subnet", resourceNsxtPolicyIPPoolStaticSubnet, resourceTypeIpAddressPoolStaticSubnet},
		{"policy_lb_pool", resourceNsxtPolicyLBPool, resourceTypeLBPool},
		{"policy_lb_server_ssl_profile", resourceNsxtPolicyLBServerSslProfile, resourceTypeLBServerSslProfile},
		{"policy_lb_service", resourceNsxtPolicyLBService, resourceTypeLBService},
		{"policy_lb_source_ip_persistence_profile", resourceNsxtPolicyLBSourceIpPersistenceProfile, resourceTypeLBSourceIpPersistenceProfile},
		{"policy_lb_tcp_monitor_profile", resourceNsxtPolicyLBTcpMonitorProfile, resourceTypeLBTcpMonitorProfile},
		{"policy_lb_udp_monitor_profile", resourceNsxtPolicyLBUdpMonitorProfile, resourceTypeLBUdpMonitorProfile},
		{"policy_lb_virtual_server", resourceNsxtPolicyLBVirtualServer, resourceTypeLBVirtualServer},
		{"policy_nat_rule", resourceNsxtPolicyNATRule, resourceTypePolicyNatRule},
		{"policy_parent_security_policy", resourceNsxtPolicyParentSecurityPolicy, resourceTypeSecurityPolicy},
		{"policy_project_ip_address_allocation", resourceNsxtPolicyProjectIpAddressAllocation, resourceTypeProjectIpAddressAllocation},
		{"policy_segment_port", resourceNsxtPolicySegmentPort, resourceTypeSegmentPort},
		{"policy_service", resourceNsxtPolicyService, resourceTypeService},
		{"policy_static_route", resourceNsxtPolicyStaticRoute, resourceTypeStaticRoutes},
		{"policy_tier1_gateway", resourceNsxtPolicyTier1Gateway, resourceTypeTier1},
		{"policy_tier1_gateway_interface", resourceNsxtPolicyTier1GatewayInterface, resourceTypeTier1Interface},
		{"policy_transit_gateway", resourceNsxtPolicyTransitGateway, resourceTypeTransitGateway},
		{"policy_transit_gateway_attachment", resourceNsxtPolicyTransitGatewayAttachment, resourceTypeTransitGatewayAttachment},
		{"vpc", resourceNsxtVpc, resourceTypeVpc},
		{"vpc_attachment", resourceNsxtVpcAttachment, resourceTypeVpcAttachment},
		{"vpc_connectivity_profile", resourceNsxtVpcConnectivityProfile, resourceTypeVpcConnectivityProfile},
		{"vpc_dhcp_v4_static_binding_config", resourceNsxtVpcSubnetDhcpV4StaticBindingConfig, resourceTypeDhcpV4StaticBindingConfig},
		{"vpc_gateway_policy", resourceNsxtVPCGatewayPolicy, resourceTypeGatewayPolicy},
		{"vpc_group", resourceNsxtVPCGroup, resourceTypeVPCGroup},
		{"vpc_ip_address_allocation", resourceNsxtVpcIpAddressAllocation, resourceTypeVpcIpAddressAllocation},
		{"vpc_nat_rule", resourceNsxtPolicyVpcNatRule, resourceTypePolicyVpcNatRule},
		{"vpc_service_profile", resourceNsxtVpcServiceProfile, resourceTypeVpcServiceProfile},
		{"vpc_static_routes", resourceNsxtVpcStaticRoutes, resourceTypeStaticRoutes},
		{"vpc_subnet", resourceNsxtVpcSubnet, resourceTypeVpcSubnet},
	} {
		t.Run(resource.name, func(t *testing.T) {
			for _, status := range []int{http.StatusServiceUnavailable, http.StatusNotFound} {
				t.Run(http.StatusText(status), func(t *testing.T) {
					res := resource.factory()
					raw := map[string]interface{}{}
					for k, v := range map[string]interface{}{
						"gateway_path": "/infra/tier-1s/gw", "segment_path": "/infra/segments/segment", "pool_path": "/infra/ip-pools/pool",
						"parent_path": "/orgs/default/projects/project/vpcs/vpc/subnets/subnet", "domain": "default", "domain_id": "default",
						"locale_service_id": "default", "gateway_id": "gw", "service_path": "/infra/lb-services/service",
						"group_path": "/infra/domains/default/groups/group", "policy_path": "/infra/domains/default/security-policies/policy",
						"security_policy_path": "/infra/domains/default/security-policies/policy",
						"context":              []interface{}{map[string]interface{}{"project_id": "project", "vpc_id": "vpc"}},
					} {
						if _, ok := res.Schema[k]; ok {
							raw[k] = v
						}
					}
					// Traditional policy resources use Local Manager; VPC/project resources need their scoped context.
					if !strings.HasPrefix(resource.name, "vpc") && resource.resourceType != resourceTypeTransitGateway && resource.resourceType != resourceTypeTransitGatewayAttachment && resource.resourceType != resourceTypeProjectIpAddressAllocation {
						delete(raw, "context")
					}
					if resource.resourceType == resourceTypeTransitGateway || resource.resourceType == resourceTypeTransitGatewayAttachment || resource.resourceType == resourceTypeProjectIpAddressAllocation {
						raw["context"] = []interface{}{map[string]interface{}{"project_id": "project"}}
					}
					switch resource.resourceType {
					case resourceTypeConnectivityPolicy:
						raw["parent_path"] = "/orgs/default/projects/project/vpcs/vpc"
					case resourceTypeTransitGatewayAttachment:
						raw["parent_path"] = "/orgs/default/projects/project/transit-gateways/gateway"
					case resourceTypeVpcAttachment:
						raw["parent_path"] = "/orgs/default/projects/project/vpcs/vpc"
					}
					d := schema.TestResourceDataRaw(t, res.Schema, raw)
					d.SetId("test-id")
					m := newIsolatedMockProviderClient(t)
					m.Host = "https://nsx.invalid"
					m.CommonConfig.CacheMode = "global"
					m.CommonConfig.MaxRetries = 0
					requests := 0
					m.PolicyHTTPClient = &http.Client{Transport: unitPolicyTransport(func(r *http.Request) (*http.Response, error) {
						requests++
						assert.Equal(t, http.MethodGet, r.Method)
						return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"error_code":500,"error_message":"backend unavailable","module_name":"unit-test"}`)), Request: r}, nil
					})}
					query := getCacheQueryKey(resource.resourceType, d, m)
					tc := gcache.getTypeCache(resource.resourceType)
					tc.data[query] = map[string]*data.StructValue{}
					err := res.Read(d, m)
					require.Positive(t, requests, "the cache miss must reach the SDK")
					if status == http.StatusNotFound && resource.resourceType != resourceTypeSegmentPort {
						require.NoError(t, err)
						assert.Empty(t, d.Id())
					} else {
						require.Error(t, err)
						assert.Equal(t, "test-id", d.Id())
					}
				})
			}
		})
	}
}
