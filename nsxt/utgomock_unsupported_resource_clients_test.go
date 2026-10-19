//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

// Client constructors return nil when the provider's scope does not support a resource.
// Replace that dependency to verify resource callbacks return an actionable error.
func unavailablePolicyClient[T any](t *testing.T, factory *func(utl.SessionContext, client.Connector) *T) {
	t.Helper()
	original := *factory
	*factory = func(utl.SessionContext, client.Connector) *T { return nil }
	t.Cleanup(func() { *factory = original })
}

func TestUnitNsxt_UnsupportedResourceClients(t *testing.T) {
	oldVersion := util.NsxVersion
	t.Cleanup(func() { util.NsxVersion = oldVersion })
	util.NsxVersion = "9.2.0"
	for _, tc := range []struct {
		name        string
		resource    func() *schema.Resource
		operation   func(*schema.ResourceData, interface{}) error
		unavailable func(*testing.T)
	}{
		{"resourceNsxtPolicyConnectivityPolicyCreate", resourceNsxtPolicyConnectivityPolicy, resourceNsxtPolicyConnectivityPolicyCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliConnectivityPoliciesClient) }},
		{"resourceNsxtPolicyConnectivityPolicyRead", resourceNsxtPolicyConnectivityPolicy, resourceNsxtPolicyConnectivityPolicyRead, func(t *testing.T) { unavailablePolicyClient(t, &cliConnectivityPoliciesClient) }},
		{"resourceNsxtPolicyConnectivityPolicyUpdate", resourceNsxtPolicyConnectivityPolicy, resourceNsxtPolicyConnectivityPolicyUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliConnectivityPoliciesClient) }},
		{"resourceNsxtPolicyConnectivityPolicyDelete", resourceNsxtPolicyConnectivityPolicy, resourceNsxtPolicyConnectivityPolicyDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliConnectivityPoliciesClient) }},
		{"resourceNsxtPolicyContextProfileRead", resourceNsxtPolicyContextProfile, resourceNsxtPolicyContextProfileRead, func(t *testing.T) { unavailablePolicyClient(t, &cliContextProfilesClient) }},
		{"resourceNsxtPolicyContextProfileDelete", resourceNsxtPolicyContextProfile, resourceNsxtPolicyContextProfileDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliContextProfilesClient) }},
		{"resourceNsxtPolicyContextProfileCustomAttributeCreate", resourceNsxtPolicyContextProfileCustomAttribute, resourceNsxtPolicyContextProfileCustomAttributeCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliDefaultClient) }},
		{"resourceNsxtPolicyDhcpRelayConfigCreate", resourceNsxtPolicyDhcpRelayConfig, resourceNsxtPolicyDhcpRelayConfigCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliDhcpRelayConfigsClient) }},
		{"resourceNsxtPolicyDhcpRelayConfigRead", resourceNsxtPolicyDhcpRelayConfig, resourceNsxtPolicyDhcpRelayConfigRead, func(t *testing.T) { unavailablePolicyClient(t, &cliDhcpRelayConfigsClient) }},
		{"resourceNsxtPolicyDhcpRelayConfigUpdate", resourceNsxtPolicyDhcpRelayConfig, resourceNsxtPolicyDhcpRelayConfigUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliDhcpRelayConfigsClient) }},
		{"resourceNsxtPolicyDhcpRelayConfigDelete", resourceNsxtPolicyDhcpRelayConfig, resourceNsxtPolicyDhcpRelayConfigDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliDhcpRelayConfigsClient) }},
		{"resourceNsxtPolicyDhcpServerCreate", resourceNsxtPolicyDhcpServer, resourceNsxtPolicyDhcpServerCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliDhcpServerConfigsClient) }},
		{"resourceNsxtPolicyDhcpServerRead", resourceNsxtPolicyDhcpServer, resourceNsxtPolicyDhcpServerRead, func(t *testing.T) { unavailablePolicyClient(t, &cliDhcpServerConfigsClient) }},
		{"resourceNsxtPolicyDhcpServerUpdate", resourceNsxtPolicyDhcpServer, resourceNsxtPolicyDhcpServerUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliDhcpServerConfigsClient) }},
		{"resourceNsxtPolicyDhcpServerDelete", resourceNsxtPolicyDhcpServer, resourceNsxtPolicyDhcpServerDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliDhcpServerConfigsClient) }},
		{"resourceNsxtPolicyDhcpV4StaticBindingRead", resourceNsxtPolicyDhcpV4StaticBinding, resourceNsxtPolicyDhcpV4StaticBindingRead, func(t *testing.T) { unavailablePolicyClient(t, &cliSegmentsDhcpStaticBindingConfigsClient) }},
		{"resourceNsxtPolicyDhcpStaticBindingDelete", resourceNsxtPolicyDhcpV4StaticBinding, resourceNsxtPolicyDhcpStaticBindingDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliSegmentsDhcpStaticBindingConfigsClient) }},
		{"resourceNsxtPolicyDhcpV6StaticBindingRead", resourceNsxtPolicyDhcpV6StaticBinding, resourceNsxtPolicyDhcpV6StaticBindingRead, func(t *testing.T) { unavailablePolicyClient(t, &cliSegmentsDhcpStaticBindingConfigsClient) }},
		{"resourceNsxtPolicyDistributedFloodProtectionProfileRead", resourceNsxtPolicyDistributedFloodProtectionProfile, resourceNsxtPolicyDistributedFloodProtectionProfileRead, func(t *testing.T) { unavailablePolicyClient(t, &cliFloodProtectionProfilesClient) }},
		{"resourceNsxtPolicyFloodProtectionProfileDelete", resourceNsxtPolicyDistributedFloodProtectionProfile, resourceNsxtPolicyFloodProtectionProfileDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliFloodProtectionProfilesClient) }},
		{"resourceNsxtPolicyDistributedFloodProtectionProfileBindingRead", resourceNsxtPolicyDistributedFloodProtectionProfileBinding, resourceNsxtPolicyDistributedFloodProtectionProfileBindingRead, func(t *testing.T) { unavailablePolicyClient(t, &cliFirewallFloodProtectionProfileBindingMapsClient) }},
		{"resourceNsxtPolicyDistributedFloodProtectionProfileBindingDelete", resourceNsxtPolicyDistributedFloodProtectionProfileBinding, resourceNsxtPolicyDistributedFloodProtectionProfileBindingDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliFirewallFloodProtectionProfileBindingMapsClient) }},
		{"resourceNsxtPolicyDistributedVlanConnectionCreate", resourceNsxtPolicyDistributedVlanConnection, resourceNsxtPolicyDistributedVlanConnectionCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliDistributedVlanConnectionsClient) }},
		{"resourceNsxtPolicyDistributedVlanConnectionRead", resourceNsxtPolicyDistributedVlanConnection, resourceNsxtPolicyDistributedVlanConnectionRead, func(t *testing.T) { unavailablePolicyClient(t, &cliDistributedVlanConnectionsClient) }},
		{"resourceNsxtPolicyDistributedVlanConnectionUpdate", resourceNsxtPolicyDistributedVlanConnection, resourceNsxtPolicyDistributedVlanConnectionUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliDistributedVlanConnectionsClient) }},
		{"resourceNsxtPolicyDistributedVlanConnectionDelete", resourceNsxtPolicyDistributedVlanConnection, resourceNsxtPolicyDistributedVlanConnectionDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliDistributedVlanConnectionsClient) }},
		{"resourceNsxtPolicyDistributedVxlanConnectionCreate", resourceNsxtPolicyDistributedVxlanConnection, resourceNsxtPolicyDistributedVxlanConnectionCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliDistributedVxlanConnectionsClient) }},
		{"resourceNsxtPolicyDistributedVxlanConnectionRead", resourceNsxtPolicyDistributedVxlanConnection, resourceNsxtPolicyDistributedVxlanConnectionRead, func(t *testing.T) { unavailablePolicyClient(t, &cliDistributedVxlanConnectionsClient) }},
		{"resourceNsxtPolicyDistributedVxlanConnectionUpdate", resourceNsxtPolicyDistributedVxlanConnection, resourceNsxtPolicyDistributedVxlanConnectionUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliDistributedVxlanConnectionsClient) }},
		{"resourceNsxtPolicyDistributedVxlanConnectionDelete", resourceNsxtPolicyDistributedVxlanConnection, resourceNsxtPolicyDistributedVxlanConnectionDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliDistributedVxlanConnectionsClient) }},
		{"resourceNsxtPolicyDNSForwarderZoneRead", resourceNsxtPolicyDNSForwarderZone, resourceNsxtPolicyDNSForwarderZoneRead, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsForwarderZonesClient) }},
		{"resourceNsxtPolicyDNSForwarderZoneDelete", resourceNsxtPolicyDNSForwarderZone, resourceNsxtPolicyDNSForwarderZoneDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsForwarderZonesClient) }},
		{"resourceNsxtPolicyDnsRecordCreate", resourceNsxtPolicyDnsRecord, resourceNsxtPolicyDnsRecordCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliProjectDnsRecordsClient) }},
		{"resourceNsxtPolicyDnsRecordRead", resourceNsxtPolicyDnsRecord, resourceNsxtPolicyDnsRecordRead, func(t *testing.T) { unavailablePolicyClient(t, &cliProjectDnsRecordsClient) }},
		{"resourceNsxtPolicyDnsRecordUpdate", resourceNsxtPolicyDnsRecord, resourceNsxtPolicyDnsRecordUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliProjectDnsRecordsClient) }},
		{"resourceNsxtPolicyDnsRecordDelete", resourceNsxtPolicyDnsRecord, resourceNsxtPolicyDnsRecordDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliProjectDnsRecordsClient) }},
		{"resourceNsxtPolicyDnsRecordAutoConfigCreate", resourceNsxtPolicyDnsRecordAutoConfig, resourceNsxtPolicyDnsRecordAutoConfigCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsAutoRecordConfigsClient) }},
		{"resourceNsxtPolicyDnsRecordAutoConfigRead", resourceNsxtPolicyDnsRecordAutoConfig, resourceNsxtPolicyDnsRecordAutoConfigRead, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsAutoRecordConfigsClient) }},
		{"resourceNsxtPolicyDnsRecordAutoConfigUpdate", resourceNsxtPolicyDnsRecordAutoConfig, resourceNsxtPolicyDnsRecordAutoConfigUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsAutoRecordConfigsClient) }},
		{"resourceNsxtPolicyDnsRecordAutoConfigDelete", resourceNsxtPolicyDnsRecordAutoConfig, resourceNsxtPolicyDnsRecordAutoConfigDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsAutoRecordConfigsClient) }},
		{"resourceNsxtPolicyDnsRuleRead", resourceNsxtPolicyDnsRule, resourceNsxtPolicyDnsRuleRead, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsRulesClient) }},
		{"resourceNsxtPolicyDnsRuleDelete", resourceNsxtPolicyDnsRule, resourceNsxtPolicyDnsRuleDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsRulesClient) }},
		{"resourceNsxtPolicyDnsServiceCreate", resourceNsxtPolicyDnsService, resourceNsxtPolicyDnsServiceCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliPolicyDnsServicesClient) }},
		{"resourceNsxtPolicyDnsServiceRead", resourceNsxtPolicyDnsService, resourceNsxtPolicyDnsServiceRead, func(t *testing.T) { unavailablePolicyClient(t, &cliPolicyDnsServicesClient) }},
		{"resourceNsxtPolicyDnsServiceUpdate", resourceNsxtPolicyDnsService, resourceNsxtPolicyDnsServiceUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliPolicyDnsServicesClient) }},
		{"resourceNsxtPolicyDnsServiceDelete", resourceNsxtPolicyDnsService, resourceNsxtPolicyDnsServiceDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliPolicyDnsServicesClient) }},
		{"resourceNsxtPolicyDnsZoneCreate", resourceNsxtPolicyDnsZone, resourceNsxtPolicyDnsZoneCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsZonesClient) }},
		{"resourceNsxtPolicyDnsZoneRead", resourceNsxtPolicyDnsZone, resourceNsxtPolicyDnsZoneRead, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsZonesClient) }},
		{"resourceNsxtPolicyDnsZoneUpdate", resourceNsxtPolicyDnsZone, resourceNsxtPolicyDnsZoneUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsZonesClient) }},
		{"resourceNsxtPolicyDnsZoneDelete", resourceNsxtPolicyDnsZone, resourceNsxtPolicyDnsZoneDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsZonesClient) }},
		{"resourceNsxtPolicyDomainRead", resourceNsxtPolicyDomain, resourceNsxtPolicyDomainRead, func(t *testing.T) { unavailablePolicyClient(t, &cliDomainDeploymentMapsClient) }},
		{"resourceNsxtPolicyFirewallExcludeListMemberCreate", resourceNsxtPolicyFirewallExcludeListMember, resourceNsxtPolicyFirewallExcludeListMemberCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliExcludeListClient) }},
		{"resourceNsxtPolicyFirewallExcludeListMemberRead", resourceNsxtPolicyFirewallExcludeListMember, resourceNsxtPolicyFirewallExcludeListMemberRead, func(t *testing.T) { unavailablePolicyClient(t, &cliExcludeListClient) }},
		{"resourceNsxtPolicyFirewallExcludeListMemberDelete", resourceNsxtPolicyFirewallExcludeListMember, resourceNsxtPolicyFirewallExcludeListMemberDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliExcludeListClient) }},
		{"resourceNsxtPolicyGatewayConnectionCreate", resourceNsxtPolicyGatewayConnection, resourceNsxtPolicyGatewayConnectionCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayConnectionsClient) }},
		{"resourceNsxtPolicyGatewayConnectionRead", resourceNsxtPolicyGatewayConnection, resourceNsxtPolicyGatewayConnectionRead, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayConnectionsClient) }},
		{"resourceNsxtPolicyGatewayConnectionUpdate", resourceNsxtPolicyGatewayConnection, resourceNsxtPolicyGatewayConnectionUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayConnectionsClient) }},
		{"resourceNsxtPolicyGatewayConnectionDelete", resourceNsxtPolicyGatewayConnection, resourceNsxtPolicyGatewayConnectionDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayConnectionsClient) }},
		{"resourceNsxtPolicyGatewayFloodProtectionProfileRead", resourceNsxtPolicyGatewayFloodProtectionProfile, resourceNsxtPolicyGatewayFloodProtectionProfileRead, func(t *testing.T) { unavailablePolicyClient(t, &cliFloodProtectionProfilesClient) }},
		{"resourceNsxtPolicyGatewayPolicyDelete", resourceNsxtPolicyGatewayPolicy, resourceNsxtPolicyGatewayPolicyDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayPoliciesClient) }},
		{"resourceNsxtPolicyGatewayPolicyRuleCreate", resourceNsxtPolicyGatewayRulePolicy, resourceNsxtPolicyGatewayPolicyRuleCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayPoliciesRulesClient) }},
		{"resourceNsxtPolicyGatewayPolicyRuleRead", resourceNsxtPolicyGatewayRulePolicy, resourceNsxtPolicyGatewayPolicyRuleRead, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayPoliciesRulesClient) }},
		{"resourceNsxtPolicyGatewayPolicyRuleUpdate", resourceNsxtPolicyGatewayRulePolicy, resourceNsxtPolicyGatewayPolicyRuleUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayPoliciesRulesClient) }},
		{"resourceNsxtPolicyGatewayPolicyRuleDelete", resourceNsxtPolicyGatewayRulePolicy, resourceNsxtPolicyGatewayPolicyRuleDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayPoliciesRulesClient) }},
		{"resourceNsxtPolicyGatewayRedistributionConfigRead", resourceNsxtPolicyGatewayRedistributionConfig, resourceNsxtPolicyGatewayRedistributionConfigRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTier0LocaleServicesClient) }},
		{"resourceNsxtPolicyGatewayRedistributionConfigDelete", resourceNsxtPolicyGatewayRedistributionConfig, resourceNsxtPolicyGatewayRedistributionConfigDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTier0LocaleServicesClient) }},
		{"resourceNsxtPolicyIdpsClusterConfigCreate", resourceNsxtPolicyIdpsClusterConfig, resourceNsxtPolicyIdpsClusterConfigCreate, func(t *testing.T) { unavailableStandalonePolicyClient(t, &cliIdsClusterConfigsClient) }},
		{"resourceNsxtPolicyIdpsClusterConfigRead", resourceNsxtPolicyIdpsClusterConfig, resourceNsxtPolicyIdpsClusterConfigRead, func(t *testing.T) { unavailableStandalonePolicyClient(t, &cliIdsClusterConfigsClient) }},
		{"resourceNsxtPolicyIdpsClusterConfigUpdate", resourceNsxtPolicyIdpsClusterConfig, resourceNsxtPolicyIdpsClusterConfigUpdate, func(t *testing.T) { unavailableStandalonePolicyClient(t, &cliIdsClusterConfigsClient) }},
		{"resourceNsxtPolicyIdpsClusterConfigDelete", resourceNsxtPolicyIdpsClusterConfig, resourceNsxtPolicyIdpsClusterConfigDelete, func(t *testing.T) { unavailableStandalonePolicyClient(t, &cliIdsClusterConfigsClient) }},
		{"resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete", resourceNsxtPolicyIntrusionServiceGatewayPolicy, resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceGatewayPoliciesClient) }},
		{"resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate", resourceNsxtPolicyIntrusionServiceGatewayPolicyRule, resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceGatewayPolicyRulesClient) }},
		{"resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead", resourceNsxtPolicyIntrusionServiceGatewayPolicyRule, resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceGatewayPolicyRulesClient) }},
		{"resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate", resourceNsxtPolicyIntrusionServiceGatewayPolicyRule, resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceGatewayPolicyRulesClient) }},
		{"resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete", resourceNsxtPolicyIntrusionServiceGatewayPolicyRule, resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceGatewayPolicyRulesClient) }},
		{"resourceNsxtPolicyIntrusionServicePolicyDelete", resourceNsxtPolicyIntrusionServicePolicy, resourceNsxtPolicyIntrusionServicePolicyDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServicePoliciesClient) }},
		{"resourceNsxtPolicyIntrusionServicePolicyRuleCreate", resourceNsxtPolicyIntrusionServicePolicyRule, resourceNsxtPolicyIntrusionServicePolicyRuleCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServicePolicyRulesClient) }},
		{"resourceNsxtPolicyIntrusionServicePolicyRuleRead", resourceNsxtPolicyIntrusionServicePolicyRule, resourceNsxtPolicyIntrusionServicePolicyRuleRead, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServicePolicyRulesClient) }},
		{"resourceNsxtPolicyIntrusionServicePolicyRuleUpdate", resourceNsxtPolicyIntrusionServicePolicyRule, resourceNsxtPolicyIntrusionServicePolicyRuleUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServicePolicyRulesClient) }},
		{"resourceNsxtPolicyIntrusionServicePolicyRuleDelete", resourceNsxtPolicyIntrusionServicePolicyRule, resourceNsxtPolicyIntrusionServicePolicyRuleDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServicePolicyRulesClient) }},
		{"resourceNsxtPolicyIntrusionServiceProfileCreate", resourceNsxtPolicyIntrusionServiceProfile, resourceNsxtPolicyIntrusionServiceProfileCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceProfilesClient) }},
		{"resourceNsxtPolicyIntrusionServiceProfileRead", resourceNsxtPolicyIntrusionServiceProfile, resourceNsxtPolicyIntrusionServiceProfileRead, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceProfilesClient) }},
		{"resourceNsxtPolicyIntrusionServiceProfileUpdate", resourceNsxtPolicyIntrusionServiceProfile, resourceNsxtPolicyIntrusionServiceProfileUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceProfilesClient) }},
		{"resourceNsxtPolicyIntrusionServiceProfileDelete", resourceNsxtPolicyIntrusionServiceProfile, resourceNsxtPolicyIntrusionServiceProfileDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceProfilesClient) }},
		{"resourceNsxtPolicyIPAddressAllocationCreate", resourceNsxtPolicyIPAddressAllocation, resourceNsxtPolicyIPAddressAllocationCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpAllocationsClient) }},
		{"resourceNsxtPolicyIPAddressAllocationRead", resourceNsxtPolicyIPAddressAllocation, resourceNsxtPolicyIPAddressAllocationRead, func(t *testing.T) { unavailablePolicyClient(t, &cliIpAllocationsClient) }},
		{"resourceNsxtPolicyIPAddressAllocationUpdate", resourceNsxtPolicyIPAddressAllocation, resourceNsxtPolicyIPAddressAllocationUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpAllocationsClient) }},
		{"resourceNsxtPolicyIPAddressAllocationDelete", resourceNsxtPolicyIPAddressAllocation, resourceNsxtPolicyIPAddressAllocationDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIpAllocationsClient) }},
		{"resourceNsxtPolicyIPBlockRead", resourceNsxtPolicyIPBlock, resourceNsxtPolicyIPBlockRead, func(t *testing.T) { unavailablePolicyClient(t, &cliIpBlocksClient) }},
		{"resourceNsxtPolicyIPBlockCreate", resourceNsxtPolicyIPBlock, resourceNsxtPolicyIPBlockCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpBlocksClient) }},
		{"resourceNsxtPolicyIPBlockUpdate", resourceNsxtPolicyIPBlock, resourceNsxtPolicyIPBlockUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpBlocksClient) }},
		{"resourceNsxtPolicyIPBlockDelete", resourceNsxtPolicyIPBlock, resourceNsxtPolicyIPBlockDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIpBlocksClient) }},
		{"resourceNsxtPolicyIPDiscoveryProfileCreate", resourceNsxtPolicyIPDiscoveryProfile, resourceNsxtPolicyIPDiscoveryProfileCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpDiscoveryProfilesClient) }},
		{"resourceNsxtPolicyIPDiscoveryProfileRead", resourceNsxtPolicyIPDiscoveryProfile, resourceNsxtPolicyIPDiscoveryProfileRead, func(t *testing.T) { unavailablePolicyClient(t, &cliIpDiscoveryProfilesClient) }},
		{"resourceNsxtPolicyIPDiscoveryProfileUpdate", resourceNsxtPolicyIPDiscoveryProfile, resourceNsxtPolicyIPDiscoveryProfileUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpDiscoveryProfilesClient) }},
		{"resourceNsxtPolicyIPDiscoveryProfileDelete", resourceNsxtPolicyIPDiscoveryProfile, resourceNsxtPolicyIPDiscoveryProfileDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIpDiscoveryProfilesClient) }},
		{"resourceNsxtPolicyIPPoolRead", resourceNsxtPolicyIPPool, resourceNsxtPolicyIPPoolRead, func(t *testing.T) { unavailablePolicyClient(t, &cliIpPoolsClient) }},
		{"resourceNsxtPolicyIPPoolCreate", resourceNsxtPolicyIPPool, resourceNsxtPolicyIPPoolCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpPoolsClient) }},
		{"resourceNsxtPolicyIPPoolUpdate", resourceNsxtPolicyIPPool, resourceNsxtPolicyIPPoolUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpPoolsClient) }},
		{"resourceNsxtPolicyIPPoolDelete", resourceNsxtPolicyIPPool, resourceNsxtPolicyIPPoolDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIpPoolsClient) }},
		{"resourceNsxtPolicyIPPoolBlockSubnetRead", resourceNsxtPolicyIPPoolBlockSubnet, resourceNsxtPolicyIPPoolBlockSubnetRead, func(t *testing.T) { unavailablePolicyClient(t, &cliIpSubnetsClient) }},
		{"resourceNsxtPolicyIPPoolBlockSubnetCreate", resourceNsxtPolicyIPPoolBlockSubnet, resourceNsxtPolicyIPPoolBlockSubnetCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpSubnetsClient) }},
		{"resourceNsxtPolicyIPPoolBlockSubnetUpdate", resourceNsxtPolicyIPPoolBlockSubnet, resourceNsxtPolicyIPPoolBlockSubnetUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpSubnetsClient) }},
		{"resourceNsxtPolicyIPPoolBlockSubnetDelete", resourceNsxtPolicyIPPoolBlockSubnet, resourceNsxtPolicyIPPoolBlockSubnetDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIpSubnetsClient) }},
		{"resourceNsxtPolicyIPPoolStaticSubnetRead", resourceNsxtPolicyIPPoolStaticSubnet, resourceNsxtPolicyIPPoolStaticSubnetRead, func(t *testing.T) { unavailablePolicyClient(t, &cliIpSubnetsClient) }},
		{"resourceNsxtPolicyIPPoolStaticSubnetCreate", resourceNsxtPolicyIPPoolStaticSubnet, resourceNsxtPolicyIPPoolStaticSubnetCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpSubnetsClient) }},
		{"resourceNsxtPolicyIPPoolStaticSubnetUpdate", resourceNsxtPolicyIPPoolStaticSubnet, resourceNsxtPolicyIPPoolStaticSubnetUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliIpSubnetsClient) }},
		{"resourceNsxtPolicyIPPoolStaticSubnetDelete", resourceNsxtPolicyIPPoolStaticSubnet, resourceNsxtPolicyIPPoolStaticSubnetDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliIpSubnetsClient) }},
		{"resourceNsxtPolicyLBPoolCreate", resourceNsxtPolicyLBPool, resourceNsxtPolicyLBPoolCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliLbPoolsClient) }},
		{"resourceNsxtPolicyLBPoolRead", resourceNsxtPolicyLBPool, resourceNsxtPolicyLBPoolRead, func(t *testing.T) { unavailablePolicyClient(t, &cliLbPoolsClient) }},
		{"resourceNsxtPolicyLBPoolUpdate", resourceNsxtPolicyLBPool, resourceNsxtPolicyLBPoolUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliLbPoolsClient) }},
		{"resourceNsxtPolicyLBPoolDelete", resourceNsxtPolicyLBPool, resourceNsxtPolicyLBPoolDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliLbPoolsClient) }},
		{"resourceNsxtPolicyLBServiceCreate", resourceNsxtPolicyLBService, resourceNsxtPolicyLBServiceCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliLbServicesClient) }},
		{"resourceNsxtPolicyLBServiceRead", resourceNsxtPolicyLBService, resourceNsxtPolicyLBServiceRead, func(t *testing.T) { unavailablePolicyClient(t, &cliLbServicesClient) }},
		{"resourceNsxtPolicyLBServiceUpdate", resourceNsxtPolicyLBService, resourceNsxtPolicyLBServiceUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliLbServicesClient) }},
		{"resourceNsxtPolicyLBServiceDelete", resourceNsxtPolicyLBService, resourceNsxtPolicyLBServiceDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliLbServicesClient) }},
		{"resourceNsxtPolicyLBVirtualServerCreate", resourceNsxtPolicyLBVirtualServer, resourceNsxtPolicyLBVirtualServerCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliLbVirtualServersClient) }},
		{"resourceNsxtPolicyLBVirtualServerRead", resourceNsxtPolicyLBVirtualServer, resourceNsxtPolicyLBVirtualServerRead, func(t *testing.T) { unavailablePolicyClient(t, &cliLbVirtualServersClient) }},
		{"resourceNsxtPolicyLBVirtualServerUpdate", resourceNsxtPolicyLBVirtualServer, resourceNsxtPolicyLBVirtualServerUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliLbVirtualServersClient) }},
		{"resourceNsxtPolicyLBVirtualServerDelete", resourceNsxtPolicyLBVirtualServer, resourceNsxtPolicyLBVirtualServerDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliLbVirtualServersClient) }},
		{"resourceNsxtPolicyMacDiscoveryProfileCreate", resourceNsxtPolicyMacDiscoveryProfile, resourceNsxtPolicyMacDiscoveryProfileCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliMacDiscoveryProfilesClient) }},
		{"resourceNsxtPolicyMacDiscoveryProfileRead", resourceNsxtPolicyMacDiscoveryProfile, resourceNsxtPolicyMacDiscoveryProfileRead, func(t *testing.T) { unavailablePolicyClient(t, &cliMacDiscoveryProfilesClient) }},
		{"resourceNsxtPolicyMacDiscoveryProfileUpdate", resourceNsxtPolicyMacDiscoveryProfile, resourceNsxtPolicyMacDiscoveryProfileUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliMacDiscoveryProfilesClient) }},
		{"resourceNsxtPolicyMacDiscoveryProfileDelete", resourceNsxtPolicyMacDiscoveryProfile, resourceNsxtPolicyMacDiscoveryProfileDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliMacDiscoveryProfilesClient) }},
		{"resourceNsxtPolicyNetworkSpanCreate", resourceNsxtPolicyNetworkSpan, resourceNsxtPolicyNetworkSpanCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliNetworkSpansClient) }},
		{"resourceNsxtPolicyNetworkSpanRead", resourceNsxtPolicyNetworkSpan, resourceNsxtPolicyNetworkSpanRead, func(t *testing.T) { unavailablePolicyClient(t, &cliNetworkSpansClient) }},
		{"resourceNsxtPolicyNetworkSpanUpdate", resourceNsxtPolicyNetworkSpan, resourceNsxtPolicyNetworkSpanUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliNetworkSpansClient) }},
		{"resourceNsxtPolicyNetworkSpanDelete", resourceNsxtPolicyNetworkSpan, resourceNsxtPolicyNetworkSpanDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliNetworkSpansClient) }},
		{"resourceNsxtPolicyPredefinedGatewayPolicyRead", resourceNsxtPolicyPredefinedGatewayPolicy, resourceNsxtPolicyPredefinedGatewayPolicyRead, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayPoliciesClient) }},
		{"resourceNsxtPolicyPredefinedSecurityPolicyRead", resourceNsxtPolicyPredefinedSecurityPolicy, resourceNsxtPolicyPredefinedSecurityPolicyRead, func(t *testing.T) { unavailablePolicyClient(t, &cliSecurityPoliciesClient) }},
		{"resourceNsxtPolicyProjectRead", resourceNsxtPolicyProject, resourceNsxtPolicyProjectRead, func(t *testing.T) { unavailablePolicyClient(t, &cliProjectsClient) }},
		{"resourceNsxtPolicyProjectDelete", resourceNsxtPolicyProject, resourceNsxtPolicyProjectDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliProjectsClient) }},
		{"resourceNsxtPolicyProjectCNACreate", resourceNsxtPolicyProjectCentralizedNetworkAttachment, resourceNsxtPolicyProjectCNACreate, func(t *testing.T) { unavailablePolicyClient(t, &cliCNAClient) }},
		{"resourceNsxtPolicyProjectCNARead", resourceNsxtPolicyProjectCentralizedNetworkAttachment, resourceNsxtPolicyProjectCNARead, func(t *testing.T) { unavailablePolicyClient(t, &cliCNAClient) }},
		{"resourceNsxtPolicyProjectCNAUpdate", resourceNsxtPolicyProjectCentralizedNetworkAttachment, resourceNsxtPolicyProjectCNAUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliCNAClient) }},
		{"resourceNsxtPolicyProjectCNADelete", resourceNsxtPolicyProjectCentralizedNetworkAttachment, resourceNsxtPolicyProjectCNADelete, func(t *testing.T) { unavailablePolicyClient(t, &cliCNAClient) }},
		{"resourceNsxtPolicyQosProfileCreate", resourceNsxtPolicyQosProfile, resourceNsxtPolicyQosProfileCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliQosProfilesClient) }},
		{"resourceNsxtPolicyQosProfileRead", resourceNsxtPolicyQosProfile, resourceNsxtPolicyQosProfileRead, func(t *testing.T) { unavailablePolicyClient(t, &cliQosProfilesClient) }},
		{"resourceNsxtPolicyQosProfileUpdate", resourceNsxtPolicyQosProfile, resourceNsxtPolicyQosProfileUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliQosProfilesClient) }},
		{"resourceNsxtPolicyQosProfileDelete", resourceNsxtPolicyQosProfile, resourceNsxtPolicyQosProfileDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliQosProfilesClient) }},
		{"resourceNsxtPolicyRouteControllerRead", resourceNsxtPolicyRouteController, resourceNsxtPolicyRouteControllerRead, func(t *testing.T) { unavailablePolicyClient(t, &cliRouteControllersClient) }},
		{"resourceNsxtPolicySecurityPolicyDelete", resourceNsxtPolicySecurityPolicy, resourceNsxtPolicySecurityPolicyDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliSecurityPoliciesClient) }},
		{"resourceNsxtPolicySecurityPolicyRuleCreate", resourceNsxtPolicySecurityPolicyRule, resourceNsxtPolicySecurityPolicyRuleCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliSecurityPoliciesRulesClient) }},
		{"resourceNsxtPolicySecurityPolicyRuleRead", resourceNsxtPolicySecurityPolicyRule, resourceNsxtPolicySecurityPolicyRuleRead, func(t *testing.T) { unavailablePolicyClient(t, &cliSecurityPoliciesRulesClient) }},
		{"resourceNsxtPolicySecurityPolicyRuleUpdate", resourceNsxtPolicySecurityPolicyRule, resourceNsxtPolicySecurityPolicyRuleUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliSecurityPoliciesRulesClient) }},
		{"resourceNsxtPolicySecurityPolicyRuleDelete", resourceNsxtPolicySecurityPolicyRule, resourceNsxtPolicySecurityPolicyRuleDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliSecurityPoliciesRulesClient) }},
		{"resourceNsxtPolicySegmentSecurityProfileRead", resourceNsxtPolicySegmentSecurityProfile, resourceNsxtPolicySegmentSecurityProfileRead, func(t *testing.T) { unavailablePolicyClient(t, &cliSegmentSecurityProfilesClient) }},
		{"resourceNsxtPolicySegmentSecurityProfileDelete", resourceNsxtPolicySegmentSecurityProfile, resourceNsxtPolicySegmentSecurityProfileDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliSegmentSecurityProfilesClient) }},
		{"resourceNsxtPolicyServiceCreate", resourceNsxtPolicyService, resourceNsxtPolicyServiceCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliServicesClient) }},
		{"resourceNsxtPolicyServiceRead", resourceNsxtPolicyService, resourceNsxtPolicyServiceRead, func(t *testing.T) { unavailablePolicyClient(t, &cliServicesClient) }},
		{"resourceNsxtPolicyServiceUpdate", resourceNsxtPolicyService, resourceNsxtPolicyServiceUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliServicesClient) }},
		{"resourceNsxtPolicyServiceDelete", resourceNsxtPolicyService, resourceNsxtPolicyServiceDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliServicesClient) }},
		{"resourceNsxtPolicyTier0GatewayHAVipConfigCreate", resourceNsxtPolicyTier0GatewayHAVipConfig, resourceNsxtPolicyTier0GatewayHAVipConfigCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliTier0LocaleServicesClient) }},
		{"resourceNsxtPolicyTier0GatewayInterfaceRead", resourceNsxtPolicyTier0GatewayInterface, resourceNsxtPolicyTier0GatewayInterfaceRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTier0InterfacesClient) }},
		{"resourceNsxtPolicyTier0GatewayInterfaceUpdate", resourceNsxtPolicyTier0GatewayInterface, resourceNsxtPolicyTier0GatewayInterfaceUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliTier0InterfacesClient) }},
		{"resourceNsxtPolicyTier0GatewayInterfaceDelete", resourceNsxtPolicyTier0GatewayInterface, resourceNsxtPolicyTier0GatewayInterfaceDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTier0InterfacesClient) }},
		{"resourceNsxtPolicyTier1GatewayRead", resourceNsxtPolicyTier1Gateway, resourceNsxtPolicyTier1GatewayRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTier1sClient) }},
		{"resourceNsxtPolicyTier1GatewayInterfaceUpdate", resourceNsxtPolicyTier1GatewayInterface, resourceNsxtPolicyTier1GatewayInterfaceUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliTier1InterfacesClient) }},
		{"resourceNsxtPolicyTier1GatewayInterfaceDelete", resourceNsxtPolicyTier1GatewayInterface, resourceNsxtPolicyTier1GatewayInterfaceDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTier1InterfacesClient) }},
		{"resourceNsxtPolicyTransitGatewayCreate", resourceNsxtPolicyTransitGateway, resourceNsxtPolicyTransitGatewayCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliOrgRootClient) }},
		{"resourceNsxtPolicyTransitGatewayDelete", resourceNsxtPolicyTransitGateway, resourceNsxtPolicyTransitGatewayDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliOrgRootClient) }},
		{"resourceNsxtPolicyTransitGatewayAttachmentCreate", resourceNsxtPolicyTransitGatewayAttachment, resourceNsxtPolicyTransitGatewayAttachmentCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayAttachmentsClient) }},
		{"resourceNsxtPolicyTransitGatewayAttachmentRead", resourceNsxtPolicyTransitGatewayAttachment, resourceNsxtPolicyTransitGatewayAttachmentRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayAttachmentsClient) }},
		{"resourceNsxtPolicyTransitGatewayAttachmentUpdate", resourceNsxtPolicyTransitGatewayAttachment, resourceNsxtPolicyTransitGatewayAttachmentUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayAttachmentsClient) }},
		{"resourceNsxtPolicyTransitGatewayAttachmentDelete", resourceNsxtPolicyTransitGatewayAttachment, resourceNsxtPolicyTransitGatewayAttachmentDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayAttachmentsClient) }},
		{"resourceNsxtPolicyTransitGatewayBfdPeerCreate", resourceNsxtPolicyTransitGatewayBfdPeer, resourceNsxtPolicyTransitGatewayBfdPeerCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWBfdPeersClient) }},
		{"resourceNsxtPolicyTransitGatewayBfdPeerRead", resourceNsxtPolicyTransitGatewayBfdPeer, resourceNsxtPolicyTransitGatewayBfdPeerRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWBfdPeersClient) }},
		{"resourceNsxtPolicyTransitGatewayBfdPeerUpdate", resourceNsxtPolicyTransitGatewayBfdPeer, resourceNsxtPolicyTransitGatewayBfdPeerUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWBfdPeersClient) }},
		{"resourceNsxtPolicyTransitGatewayBfdPeerDelete", resourceNsxtPolicyTransitGatewayBfdPeer, resourceNsxtPolicyTransitGatewayBfdPeerDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWBfdPeersClient) }},
		{"resourceNsxtPolicyTransitGatewayBgpNeighborCreate", resourceNsxtPolicyTransitGatewayBgpNeighbor, resourceNsxtPolicyTransitGatewayBgpNeighborCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWBgpNeighborsClient) }},
		{"resourceNsxtPolicyTransitGatewayBgpNeighborRead", resourceNsxtPolicyTransitGatewayBgpNeighbor, resourceNsxtPolicyTransitGatewayBgpNeighborRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWBgpNeighborsClient) }},
		{"resourceNsxtPolicyTransitGatewayBgpNeighborUpdate", resourceNsxtPolicyTransitGatewayBgpNeighbor, resourceNsxtPolicyTransitGatewayBgpNeighborUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWBgpNeighborsClient) }},
		{"resourceNsxtPolicyTransitGatewayBgpNeighborDelete", resourceNsxtPolicyTransitGatewayBgpNeighbor, resourceNsxtPolicyTransitGatewayBgpNeighborDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWBgpNeighborsClient) }},
		{"resourceNsxtPolicyTransitGatewayCommunityListCreate", resourceNsxtPolicyTransitGatewayCommunityList, resourceNsxtPolicyTransitGatewayCommunityListCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWCommunityListsClient) }},
		{"resourceNsxtPolicyTransitGatewayCommunityListRead", resourceNsxtPolicyTransitGatewayCommunityList, resourceNsxtPolicyTransitGatewayCommunityListRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWCommunityListsClient) }},
		{"resourceNsxtPolicyTransitGatewayCommunityListUpdate", resourceNsxtPolicyTransitGatewayCommunityList, resourceNsxtPolicyTransitGatewayCommunityListUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWCommunityListsClient) }},
		{"resourceNsxtPolicyTransitGatewayCommunityListDelete", resourceNsxtPolicyTransitGatewayCommunityList, resourceNsxtPolicyTransitGatewayCommunityListDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWCommunityListsClient) }},
		{"resourceNsxtPolicyTGWIPSecVpnLocalEndpointRead", resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint, resourceNsxtPolicyTGWIPSecVpnLocalEndpointRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayIpsecVpnLocalEndpointsClient) }},
		{"resourceNsxtPolicyTGWIPSecVpnServicesCreate", resourceNsxtPolicyTransitGatewayIpsecVpnService, resourceNsxtPolicyTGWIPSecVpnServicesCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayIpsecVpnServicesClient) }},
		{"resourceNsxtPolicyTGWIPSecVpnServicesRead", resourceNsxtPolicyTransitGatewayIpsecVpnService, resourceNsxtPolicyTGWIPSecVpnServicesRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayIpsecVpnServicesClient) }},
		{"resourceNsxtPolicyTGWIPSecVpnServicesUpdate", resourceNsxtPolicyTransitGatewayIpsecVpnService, resourceNsxtPolicyTGWIPSecVpnServicesUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayIpsecVpnServicesClient) }},
		{"resourceNsxtPolicyTGWIPSecVpnServicesDelete", resourceNsxtPolicyTransitGatewayIpsecVpnService, resourceNsxtPolicyTGWIPSecVpnServicesDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayIpsecVpnServicesClient) }},
		{"resourceNsxtPolicyTGWIPSecVpnSessionRead", resourceNsxtPolicyTransitGatewayIPSecVpnSession, resourceNsxtPolicyTGWIPSecVpnSessionRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayIpsecVpnSessionsClient) }},
		{"resourceNsxtPolicyTransitGatewayNatRuleRead", resourceNsxtPolicyTransitGatewayNatRule, resourceNsxtPolicyTransitGatewayNatRuleRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayNatRulesClient) }},
		{"resourceNsxtPolicyTransitGatewayPrefixListCreate", resourceNsxtPolicyTransitGatewayPrefixList, resourceNsxtPolicyTransitGatewayPrefixListCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWPrefixListsClient) }},
		{"resourceNsxtPolicyTransitGatewayPrefixListRead", resourceNsxtPolicyTransitGatewayPrefixList, resourceNsxtPolicyTransitGatewayPrefixListRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWPrefixListsClient) }},
		{"resourceNsxtPolicyTransitGatewayPrefixListUpdate", resourceNsxtPolicyTransitGatewayPrefixList, resourceNsxtPolicyTransitGatewayPrefixListUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWPrefixListsClient) }},
		{"resourceNsxtPolicyTransitGatewayPrefixListDelete", resourceNsxtPolicyTransitGatewayPrefixList, resourceNsxtPolicyTransitGatewayPrefixListDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWPrefixListsClient) }},
		{"resourceNsxtPolicyTransitGatewayRouteMapCreate", resourceNsxtPolicyTransitGatewayRouteMap, resourceNsxtPolicyTransitGatewayRouteMapCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWRouteMapsClient) }},
		{"resourceNsxtPolicyTransitGatewayRouteMapRead", resourceNsxtPolicyTransitGatewayRouteMap, resourceNsxtPolicyTransitGatewayRouteMapRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWRouteMapsClient) }},
		{"resourceNsxtPolicyTransitGatewayRouteMapUpdate", resourceNsxtPolicyTransitGatewayRouteMap, resourceNsxtPolicyTransitGatewayRouteMapUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWRouteMapsClient) }},
		{"resourceNsxtPolicyTransitGatewayRouteMapDelete", resourceNsxtPolicyTransitGatewayRouteMap, resourceNsxtPolicyTransitGatewayRouteMapDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWRouteMapsClient) }},
		{"resourceNsxtPolicyTransitGatewayStaticRouteCreate", resourceNsxtPolicyTransitGatewayStaticRoute, resourceNsxtPolicyTransitGatewayStaticRouteCreate, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayStaticRoutesClient) }},
		{"resourceNsxtPolicyTransitGatewayStaticRouteRead", resourceNsxtPolicyTransitGatewayStaticRoute, resourceNsxtPolicyTransitGatewayStaticRouteRead, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayStaticRoutesClient) }},
		{"resourceNsxtPolicyTransitGatewayStaticRouteUpdate", resourceNsxtPolicyTransitGatewayStaticRoute, resourceNsxtPolicyTransitGatewayStaticRouteUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayStaticRoutesClient) }},
		{"resourceNsxtPolicyTransitGatewayStaticRouteDelete", resourceNsxtPolicyTransitGatewayStaticRoute, resourceNsxtPolicyTransitGatewayStaticRouteDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayStaticRoutesClient) }},
		{"resourceNsxtPolicyVirtualNetworkApplianceDelete", resourceNsxtPolicyVirtualNetworkAppliance, resourceNsxtPolicyVirtualNetworkApplianceDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliVNACRUDClient) }},
		{"resourceNsxtPolicyVirtualNetworkApplianceClusterRead", resourceNsxtPolicyVirtualNetworkApplianceCluster, resourceNsxtPolicyVirtualNetworkApplianceClusterRead, func(t *testing.T) { unavailablePolicyClient(t, &cliVNAClustersClient) }},
		{"resourceNsxtPolicyVirtualNetworkApplianceClusterUpdate", resourceNsxtPolicyVirtualNetworkApplianceCluster, resourceNsxtPolicyVirtualNetworkApplianceClusterUpdate, func(t *testing.T) { unavailablePolicyClient(t, &cliVNAClustersClient) }},
		{"resourceNsxtPolicyVirtualNetworkApplianceClusterDelete", resourceNsxtPolicyVirtualNetworkApplianceCluster, resourceNsxtPolicyVirtualNetworkApplianceClusterDelete, func(t *testing.T) { unavailablePolicyClient(t, &cliVNAClustersClient) }},
		{"resourceNsxtVPCGroupRead", resourceNsxtVPCGroup, resourceNsxtVPCGroupRead, func(t *testing.T) { unavailablePolicyClient(t, &cliGroupsClient) }},
		{"resourceNsxtPolicyBgpNeighborExists", resourceNsxtPolicyBgpNeighbor, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyBgpNeighborExists("id", "id", "id", false, getPolicyConnector(m), d, m)
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliBgpNeighborsClient) }},
		{"resourceNsxtPolicyConnectivityPolicyExists", resourceNsxtPolicyConnectivityPolicy, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyConnectivityPolicyExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliConnectivityPoliciesClient) }},
		{"resourceNsxtPolicyContextProfileExists", resourceNsxtPolicyContextProfile, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyContextProfileExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliContextProfilesClient) }},
		{"resourceNsxtPolicyContextProfileCustomAttributeExists", resourceNsxtPolicyContextProfileCustomAttribute, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyContextProfileCustomAttributeExists(utl.SessionContext{ClientType: utl.Local}, "DOMAIN_NAME~example.com", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliDefaultClient) }},
		{"resourceNsxtPolicyDhcpRelayConfigExists", resourceNsxtPolicyDhcpRelayConfig, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDhcpRelayConfigExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliDhcpRelayConfigsClient) }},
		{"resourceNsxtPolicyDhcpServerExists", resourceNsxtPolicyDhcpServer, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDhcpServerExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliDhcpServerConfigsClient) }},
		{"resourceNsxtPolicyFloodProtectionProfileExists", resourceNsxtPolicyDistributedFloodProtectionProfile, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyFloodProtectionProfileExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliFloodProtectionProfilesClient) }},
		{"resourceNsxtPolicyDistributedFloodProtectionProfileBindingExists", resourceNsxtPolicyDistributedFloodProtectionProfileBinding, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingExists(utl.SessionContext{ClientType: utl.Local}, getPolicyConnector(m), "/orgs/default/projects/project/transit-gateways/gateway", "id")
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliFirewallFloodProtectionProfileBindingMapsClient) }},
		{"resourceNsxtPolicyDistributedVlanConnectionExists", resourceNsxtPolicyDistributedVlanConnection, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDistributedVlanConnectionExists("id", getPolicyConnector(m), false)
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliDistributedVlanConnectionsClient) }},
		{"resourceNsxtPolicyDistributedVxlanConnectionExists", resourceNsxtPolicyDistributedVxlanConnection, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDistributedVxlanConnectionExists("id", getPolicyConnector(m), false)
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliDistributedVxlanConnectionsClient) }},
		{"resourceNsxtPolicyDNSForwarderZoneExists", resourceNsxtPolicyDNSForwarderZone, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDNSForwarderZoneExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsForwarderZonesClient) }},
		{"resourceNsxtPolicyDnsRecordExists", resourceNsxtPolicyDnsRecord, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDnsRecordExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliProjectDnsRecordsClient) }},
		{"resourceNsxtPolicyDnsRecordAutoConfigExists", resourceNsxtPolicyDnsRecordAutoConfig, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDnsRecordAutoConfigExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsAutoRecordConfigsClient) }},
		{"resourceNsxtPolicyDnsRuleExists", resourceNsxtPolicyDnsRule, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDnsRuleExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsRulesClient) }},
		{"resourceNsxtPolicyDnsServiceExists", resourceNsxtPolicyDnsService, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDnsServiceExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliPolicyDnsServicesClient) }},
		{"resourceNsxtPolicyDnsZoneExists", resourceNsxtPolicyDnsZone, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyDnsZoneExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliDnsZonesClient) }},
		{"resourceNsxtPolicyGatewayConnectionExists", resourceNsxtPolicyGatewayConnection, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyGatewayConnectionExists("id", getPolicyConnector(m), false)
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayConnectionsClient) }},
		{"resourceNsxtPolicyGatewayPolicyRuleExists", resourceNsxtPolicyGatewayRulePolicy, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyGatewayPolicyRuleExists(utl.SessionContext{ClientType: utl.Local}, "id", "/orgs/default/projects/project/transit-gateways/gateway", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliGatewayPoliciesRulesClient) }},
		{"resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleExists", resourceNsxtPolicyIntrusionServiceGatewayPolicyRule, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleExists(utl.SessionContext{ClientType: utl.Local}, "id", "/orgs/default/projects/project/transit-gateways/gateway", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceGatewayPolicyRulesClient) }},
		{"resourceNsxtPolicyIntrusionServicePolicyRuleExists", resourceNsxtPolicyIntrusionServicePolicyRule, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyIntrusionServicePolicyRuleExists(utl.SessionContext{ClientType: utl.Local}, "id", "/orgs/default/projects/project/transit-gateways/gateway", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServicePolicyRulesClient) }},
		{"resourceNsxtPolicyIntrusionServiceProfileExists", resourceNsxtPolicyIntrusionServiceProfile, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyIntrusionServiceProfileExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliIntrusionServiceProfilesClient) }},
		{"resourceNsxtPolicyIPAddressAllocationExists", resourceNsxtPolicyIPAddressAllocation, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyIPAddressAllocationExists(utl.SessionContext{ClientType: utl.Local}, "id", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliIpAllocationsClient) }},
		{"resourceNsxtPolicyIPBlockExists", resourceNsxtPolicyIPBlock, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyIPBlockExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliIpBlocksClient) }},
		{"resourceNsxtPolicyIPDiscoveryProfileExists", resourceNsxtPolicyIPDiscoveryProfile, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyIPDiscoveryProfileExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliIpDiscoveryProfilesClient) }},
		{"resourceNsxtPolicyIPPoolExists", resourceNsxtPolicyIPPool, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyIPPoolExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliIpPoolsClient) }},
		{"resourceNsxtPolicyLBVirtualServerExists", resourceNsxtPolicyLBVirtualServer, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyLBVirtualServerExists("id", getPolicyConnector(m), false)
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliLbVirtualServersClient) }},
		{"resourceNsxtPolicyMacDiscoveryProfileExists", resourceNsxtPolicyMacDiscoveryProfile, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyMacDiscoveryProfileExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliMacDiscoveryProfilesClient) }},
		{"resourceNsxtPolicyNetworkSpanExists", resourceNsxtPolicyNetworkSpan, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyNetworkSpanExists("id", getPolicyConnector(m), false)
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliNetworkSpansClient) }},
		{"resourceNsxtPolicyProjectExists", resourceNsxtPolicyProject, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyProjectExists("id", getPolicyConnector(m), false)
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliProjectsClient) }},
		{"resourceNsxtPolicyProjectCNAExists", resourceNsxtPolicyProjectCentralizedNetworkAttachment, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyProjectCNAExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliCNAClient) }},
		{"resourceNsxtPolicyQosProfileExists", resourceNsxtPolicyQosProfile, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyQosProfileExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliQosProfilesClient) }},
		{"resourceNsxtPolicyRouteControllerExists", resourceNsxtPolicyRouteController, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyRouteControllerExists("id", getPolicyConnector(m), false)
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliRouteControllersClient) }},
		{"resourceNsxtPolicyRouteControllerBgpNeighborExists", resourceNsxtPolicyRouteControllerBgpNeighbor, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyRouteControllerBgpNeighborExists(utl.SessionContext{ClientType: utl.Local}, "/infra/route-controllers/controller/bgp", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliRCBgpNeighborClient) }},
		{"resourceNsxtPolicyRouteControllerInterfaceExists", resourceNsxtPolicyRouteControllerInterface, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyRouteControllerInterfaceExists(utl.SessionContext{ClientType: utl.Local}, "/infra/route-controllers/controller", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliRCInterfaceClient) }},
		{"resourceNsxtPolicySecurityPolicyRuleExists", resourceNsxtPolicySecurityPolicyRule, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicySecurityPolicyRuleExists(utl.SessionContext{ClientType: utl.Local}, "id", "/orgs/default/projects/project/transit-gateways/gateway", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliSecurityPoliciesRulesClient) }},
		{"resourceNsxtPolicySegmentSecurityProfileExists", resourceNsxtPolicySegmentSecurityProfile, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicySegmentSecurityProfileExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliSegmentSecurityProfilesClient) }},
		{"resourceNsxtPolicyServiceExists", resourceNsxtPolicyService, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyServiceExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliServicesClient) }},
		{"resourceNsxtPolicyTier1GatewayExists", resourceNsxtPolicyTier1Gateway, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTier1GatewayExists(utl.SessionContext{ClientType: utl.Local}, "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTier1sClient) }},
		{"resourceNsxtPolicyTransitGatewayAttachmentExists", resourceNsxtPolicyTransitGatewayAttachment, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTransitGatewayAttachmentExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayAttachmentsClient) }},
		{"resourceNsxtPolicyTransitGatewayBfdPeerExists", resourceNsxtPolicyTransitGatewayBfdPeer, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTransitGatewayBfdPeerExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWBfdPeersClient) }},
		{"resourceNsxtPolicyTransitGatewayBgpNeighborExists", resourceNsxtPolicyTransitGatewayBgpNeighbor, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTransitGatewayBgpNeighborExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWBgpNeighborsClient) }},
		{"resourceNsxtPolicyTransitGatewayCommunityListExists", resourceNsxtPolicyTransitGatewayCommunityList, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTransitGatewayCommunityListExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWCommunityListsClient) }},
		{"resourceNsxtPolicyTGWIPSecVpnLocalEndpointExists", resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway/ipsec-vpn-services/service", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayIpsecVpnLocalEndpointsClient) }},
		{"resourceNsxtPolicyTGWIPSecVpnServicesExists", resourceNsxtPolicyTransitGatewayIpsecVpnService, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTGWIPSecVpnServicesExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayIpsecVpnServicesClient) }},
		{"resourceNsxtPolicyTGWIPSecVpnSessionExists", resourceNsxtPolicyTransitGatewayIPSecVpnSession, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTGWIPSecVpnSessionExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway/ipsec-vpn-services/service", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayIpsecVpnSessionsClient) }},
		{"resourceNsxtPolicyTransitGatewayNatRuleExists", resourceNsxtPolicyTransitGatewayNatRule, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTransitGatewayNatRuleExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway/nat/default", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayNatRulesClient) }},
		{"resourceNsxtPolicyTransitGatewayPrefixListExists", resourceNsxtPolicyTransitGatewayPrefixList, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTransitGatewayPrefixListExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWPrefixListsClient) }},
		{"resourceNsxtPolicyTransitGatewayRouteMapExists", resourceNsxtPolicyTransitGatewayRouteMap, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTransitGatewayRouteMapExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTGWRouteMapsClient) }},
		{"resourceNsxtPolicyTransitGatewayStaticRouteExists", resourceNsxtPolicyTransitGatewayStaticRoute, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyTransitGatewayStaticRouteExists(utl.SessionContext{ClientType: utl.Local}, "/orgs/default/projects/project/transit-gateways/gateway", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliTransitGatewayStaticRoutesClient) }},
		{"resourceNsxtPolicyVirtualNetworkApplianceExists", resourceNsxtPolicyVirtualNetworkAppliance, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyVirtualNetworkApplianceExists("/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/cluster", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliVNACRUDClient) }},
		{"resourceNsxtPolicyVirtualNetworkApplianceClusterExists", resourceNsxtPolicyVirtualNetworkApplianceCluster, func(d *schema.ResourceData, m interface{}) error {
			_, err := resourceNsxtPolicyVirtualNetworkApplianceClusterExists("id", "id", "id", getPolicyConnector(m))
			return err
		}, func(t *testing.T) { unavailablePolicyClient(t, &cliSitesClient) }},
		{"dataSourceNsxtPolicyGatewayInterfaceRealizationRead", dataSourceNsxtPolicyGatewayInterfaceRealization, dataSourceNsxtPolicyGatewayInterfaceRealizationRead, func(t *testing.T) { unavailablePolicyClient(t, &cliRealizedEntitiesClient) }},
		{"dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembersRead", dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers, dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembersRead, func(t *testing.T) { unavailablePolicyClient(t, &cliGroupBareMetalServerInterfaceMembersClient) }},
		{"dataSourceNsxtPolicyGroupBareMetalServerMembersRead", dataSourceNsxtPolicyGroupBareMetalServerMembers, dataSourceNsxtPolicyGroupBareMetalServerMembersRead, func(t *testing.T) { unavailablePolicyClient(t, &cliGroupBareMetalServerMembersClient) }},
		{"dataSourceNsxtPolicyIPBlockRead", dataSourceNsxtPolicyIPBlock, dataSourceNsxtPolicyIPBlockRead, func(t *testing.T) { unavailablePolicyClient(t, &cliIpBlocksClient) }},
		{"dataSourceNsxtPolicyProjectRead", dataSourceNsxtPolicyProject, dataSourceNsxtPolicyProjectRead, func(t *testing.T) { unavailablePolicyClient(t, &cliProjectsClient) }},
		{"dataSourceNsxtPolicyRoleBindingRead", dataSourceNsxtPolicyRoleBinding, dataSourceNsxtPolicyRoleBindingRead, func(t *testing.T) { unavailablePolicyClient(t, &cliRoleBindingsClient) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.unavailable(t)
			res := tc.resource()
			raw := map[string]interface{}{}
			for k, v := range map[string]interface{}{
				"nsx_id": "test-id", "gateway_path": "/orgs/default/projects/project/transit-gateways/gateway", "segment_path": "/infra/segments/segment", "pool_path": "/infra/ip-pools/pool",
				"parent_path": "/orgs/default/projects/project/vpcs/vpc", "domain": "default", "domain_id": "default", "group_path": "/infra/domains/default/groups/group",
				"locale_service_id": "default", "gateway_id": "gateway", "service_path": "/infra/tier-1s/gateway/ipsec-vpn-services/vpn",
				"security_policy_path": "/infra/domains/default/security-policies/policy", "site_path": "/infra/sites/default", "enforcement_point": "default",
			} {
				if _, ok := res.Schema[k]; ok {
					raw[k] = v
				}
			}
			if strings.HasPrefix(tc.name, "dataSource") {
				if _, ok := res.Schema["id"]; ok {
					raw["id"] = "test-id"
				}
				if _, ok := res.Schema["display_name"]; ok {
					raw["display_name"] = "test"
				}
			}
			if strings.HasSuffix(tc.name, "Create") {
				delete(raw, "nsx_id")
			}
			d := schema.TestResourceDataRaw(t, res.Schema, raw)
			d.SetId("test-id")
			err := tc.operation(d, newIsolatedMockProviderClient(t))
			require.Error(t, err)
			require.True(t, strings.Contains(err.Error(), "not supported") || strings.Contains(err.Error(), "unsupported client"), "unexpected error: %v", err)
		})
	}
}

func unavailableStandalonePolicyClient[T any](t *testing.T, factory *func(client.Connector) T) {
	t.Helper()
	original := *factory
	*factory = func(client.Connector) T { var zero T; return zero }
	t.Cleanup(func() { *factory = original })
}
