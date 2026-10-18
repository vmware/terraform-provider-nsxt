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

func TestUnitNsxt_ResourceRejectsMissingIDAndInvalidParent(t *testing.T) {
	oldVersion := util.NsxVersion
	t.Cleanup(func() { util.NsxVersion = oldVersion })
	util.NsxVersion = "9.2.0"
	for _, tc := range []struct {
		name         string
		resource     func() *schema.Resource
		operation    func(*schema.ResourceData, interface{}) error
		field, value string
		missingID    bool
	}{
		{"resourceNsxtPolicyConnectivityPolicyUpdate_invalid_parent_path", resourceNsxtPolicyConnectivityPolicy, resourceNsxtPolicyConnectivityPolicyUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyConnectivityPolicyDelete_invalid_parent_path", resourceNsxtPolicyConnectivityPolicy, resourceNsxtPolicyConnectivityPolicyDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyDhcpV4StaticBindingUpdate_missing_id", resourceNsxtPolicyDhcpV4StaticBinding, resourceNsxtPolicyDhcpV4StaticBindingUpdate, "", "", true},
		{"resourceNsxtPolicyDhcpV6StaticBindingUpdate_missing_id", resourceNsxtPolicyDhcpV6StaticBinding, resourceNsxtPolicyDhcpV6StaticBindingUpdate, "", "", true},
		{"resourceNsxtPolicyDnsRecordRead_missing_id", resourceNsxtPolicyDnsRecord, resourceNsxtPolicyDnsRecordRead, "", "", true},
		{"resourceNsxtPolicyDnsRecordUpdate_missing_id", resourceNsxtPolicyDnsRecord, resourceNsxtPolicyDnsRecordUpdate, "", "", true},
		{"resourceNsxtPolicyDnsRecordDelete_missing_id", resourceNsxtPolicyDnsRecord, resourceNsxtPolicyDnsRecordDelete, "", "", true},
		{"resourceNsxtPolicyDnsRecordAutoConfigRead_missing_id", resourceNsxtPolicyDnsRecordAutoConfig, resourceNsxtPolicyDnsRecordAutoConfigRead, "", "", true},
		{"resourceNsxtPolicyDnsRecordAutoConfigUpdate_missing_id", resourceNsxtPolicyDnsRecordAutoConfig, resourceNsxtPolicyDnsRecordAutoConfigUpdate, "", "", true},
		{"resourceNsxtPolicyDnsRecordAutoConfigDelete_missing_id", resourceNsxtPolicyDnsRecordAutoConfig, resourceNsxtPolicyDnsRecordAutoConfigDelete, "", "", true},
		{"resourceNsxtPolicyDnsRuleRead_missing_id", resourceNsxtPolicyDnsRule, resourceNsxtPolicyDnsRuleRead, "", "", true},
		{"resourceNsxtPolicyDnsRuleUpdate_missing_id", resourceNsxtPolicyDnsRule, resourceNsxtPolicyDnsRuleUpdate, "", "", true},
		{"resourceNsxtPolicyDnsRuleDelete_missing_id", resourceNsxtPolicyDnsRule, resourceNsxtPolicyDnsRuleDelete, "", "", true},
		{"resourceNsxtPolicyDnsServiceRead_missing_id", resourceNsxtPolicyDnsService, resourceNsxtPolicyDnsServiceRead, "", "", true},
		{"resourceNsxtPolicyDnsServiceUpdate_missing_id", resourceNsxtPolicyDnsService, resourceNsxtPolicyDnsServiceUpdate, "", "", true},
		{"resourceNsxtPolicyDnsServiceDelete_missing_id", resourceNsxtPolicyDnsService, resourceNsxtPolicyDnsServiceDelete, "", "", true},
		{"resourceNsxtPolicyDnsZoneRead_missing_id", resourceNsxtPolicyDnsZone, resourceNsxtPolicyDnsZoneRead, "", "", true},
		{"resourceNsxtPolicyDnsZoneUpdate_missing_id", resourceNsxtPolicyDnsZone, resourceNsxtPolicyDnsZoneUpdate, "", "", true},
		{"resourceNsxtPolicyDnsZoneDelete_missing_id", resourceNsxtPolicyDnsZone, resourceNsxtPolicyDnsZoneDelete, "", "", true},
		{"resourceNsxtPolicyEvpnConfigRead_invalid_gateway_path", resourceNsxtPolicyEvpnConfig, resourceNsxtPolicyEvpnConfigRead, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyEvpnConfigUpdate_invalid_gateway_path", resourceNsxtPolicyEvpnConfig, resourceNsxtPolicyEvpnConfigUpdate, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyEvpnConfigDelete_invalid_gateway_path", resourceNsxtPolicyEvpnConfig, resourceNsxtPolicyEvpnConfigDelete, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyEvpnTunnelEndpointUpdate_missing_id", resourceNsxtPolicyEvpnTunnelEndpoint, resourceNsxtPolicyEvpnTunnelEndpointUpdate, "", "", true},
		{"resourceNsxtPolicyEvpnTunnelEndpointDelete_missing_id", resourceNsxtPolicyEvpnTunnelEndpoint, resourceNsxtPolicyEvpnTunnelEndpointDelete, "", "", true},
		{"resourceNsxtPolicyGatewayDNSForwarderRead_invalid_gateway_path", resourceNsxtPolicyGatewayDNSForwarder, resourceNsxtPolicyGatewayDNSForwarderRead, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyGatewayDNSForwarderCreate_invalid_gateway_path", resourceNsxtPolicyGatewayDNSForwarder, resourceNsxtPolicyGatewayDNSForwarderCreate, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyGatewayDNSForwarderUpdate_invalid_gateway_path", resourceNsxtPolicyGatewayDNSForwarder, resourceNsxtPolicyGatewayDNSForwarderUpdate, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyGatewayDNSForwarderDelete_invalid_gateway_path", resourceNsxtPolicyGatewayDNSForwarder, resourceNsxtPolicyGatewayDNSForwarderDelete, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyGatewayPrefixListDelete_invalid_gateway_path", resourceNsxtPolicyGatewayPrefixList, resourceNsxtPolicyGatewayPrefixListDelete, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyGatewayPrefixListRead_invalid_gateway_path", resourceNsxtPolicyGatewayPrefixList, resourceNsxtPolicyGatewayPrefixListRead, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyGatewayPrefixListCreate_invalid_gateway_path", resourceNsxtPolicyGatewayPrefixList, resourceNsxtPolicyGatewayPrefixListCreate, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyGatewayPrefixListUpdate_invalid_gateway_path", resourceNsxtPolicyGatewayPrefixList, resourceNsxtPolicyGatewayPrefixListUpdate, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyGatewayRedistributionConfigUpdate_missing_id", resourceNsxtPolicyGatewayRedistributionConfig, resourceNsxtPolicyGatewayRedistributionConfigUpdate, "", "", true},
		{"resourceNsxtPolicyGatewayRedistributionConfigDelete_missing_id", resourceNsxtPolicyGatewayRedistributionConfig, resourceNsxtPolicyGatewayRedistributionConfigDelete, "", "", true},
		{"resourceNsxtPolicyGatewaySecurityConfigRead_missing_id", resourceNsxtPolicyGatewaySecurityConfig, resourceNsxtPolicyGatewaySecurityConfigRead, "", "", true},
		{"resourceNsxtPolicyGatewaySecurityConfigUpdate_missing_id", resourceNsxtPolicyGatewaySecurityConfig, resourceNsxtPolicyGatewaySecurityConfigUpdate, "", "", true},
		{"resourceNsxtPolicyGatewaySecurityConfigDelete_missing_id", resourceNsxtPolicyGatewaySecurityConfig, resourceNsxtPolicyGatewaySecurityConfigDelete, "", "", true},
		{"resourceNsxtPolicyIdpsSettingsUpdate_missing_id", resourceNsxtPolicyIdpsSettings, resourceNsxtPolicyIdpsSettingsUpdate, "", "", true},
		{"resourceNsxtPolicyIdpsSettingsDelete_missing_id", resourceNsxtPolicyIdpsSettings, resourceNsxtPolicyIdpsSettingsDelete, "", "", true},
		{"resourceNsxtPolicyIdpsSignatureVersionUpdate_missing_id", resourceNsxtPolicyIdpsSignatureVersion, resourceNsxtPolicyIdpsSignatureVersionUpdate, "", "", true},
		{"resourceNsxtPolicyIdpsSignatureVersionDelete_missing_id", resourceNsxtPolicyIdpsSignatureVersion, resourceNsxtPolicyIdpsSignatureVersionDelete, "", "", true},
		{"resourceNsxtPolicyIPSecVpnSessionUpdate_missing_id", resourceNsxtPolicyIPSecVpnSession, resourceNsxtPolicyIPSecVpnSessionUpdate, "", "", true},
		{"resourceNsxtPolicyRouteControllerBgpNeighborCreate_invalid_parent_path", resourceNsxtPolicyRouteControllerBgpNeighbor, resourceNsxtPolicyRouteControllerBgpNeighborCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyRouteControllerBgpNeighborRead_invalid_parent_path", resourceNsxtPolicyRouteControllerBgpNeighbor, resourceNsxtPolicyRouteControllerBgpNeighborRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyRouteControllerBgpNeighborUpdate_invalid_parent_path", resourceNsxtPolicyRouteControllerBgpNeighbor, resourceNsxtPolicyRouteControllerBgpNeighborUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyRouteControllerBgpNeighborDelete_invalid_parent_path", resourceNsxtPolicyRouteControllerBgpNeighbor, resourceNsxtPolicyRouteControllerBgpNeighborDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyRouteControllerInterfaceCreate_invalid_parent_path", resourceNsxtPolicyRouteControllerInterface, resourceNsxtPolicyRouteControllerInterfaceCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyRouteControllerInterfaceRead_invalid_parent_path", resourceNsxtPolicyRouteControllerInterface, resourceNsxtPolicyRouteControllerInterfaceRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyRouteControllerInterfaceUpdate_invalid_parent_path", resourceNsxtPolicyRouteControllerInterface, resourceNsxtPolicyRouteControllerInterfaceUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyRouteControllerInterfaceDelete_invalid_parent_path", resourceNsxtPolicyRouteControllerInterface, resourceNsxtPolicyRouteControllerInterfaceDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyStaticRouteCreate_invalid_gateway_path", resourceNsxtPolicyStaticRoute, resourceNsxtPolicyStaticRouteCreate, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyStaticRouteRead_invalid_gateway_path", resourceNsxtPolicyStaticRoute, resourceNsxtPolicyStaticRouteRead, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyStaticRouteUpdate_invalid_gateway_path", resourceNsxtPolicyStaticRoute, resourceNsxtPolicyStaticRouteUpdate, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyStaticRouteDelete_invalid_gateway_path", resourceNsxtPolicyStaticRoute, resourceNsxtPolicyStaticRouteDelete, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyTier0InterVRFRoutingCreate_invalid_gateway_path", resourceNsxtPolicyTier0InterVRFRouting, resourceNsxtPolicyTier0InterVRFRoutingCreate, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyTier0InterVRFRoutingRead_invalid_gateway_path", resourceNsxtPolicyTier0InterVRFRouting, resourceNsxtPolicyTier0InterVRFRoutingRead, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyTier0InterVRFRoutingUpdate_invalid_gateway_path", resourceNsxtPolicyTier0InterVRFRouting, resourceNsxtPolicyTier0InterVRFRoutingUpdate, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyTier0InterVRFRoutingDelete_invalid_gateway_path", resourceNsxtPolicyTier0InterVRFRouting, resourceNsxtPolicyTier0InterVRFRoutingDelete, "gateway_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayAttachmentCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayAttachment, resourceNsxtPolicyTransitGatewayAttachmentCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayAttachmentRead_invalid_parent_path", resourceNsxtPolicyTransitGatewayAttachment, resourceNsxtPolicyTransitGatewayAttachmentRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayAttachmentUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayAttachment, resourceNsxtPolicyTransitGatewayAttachmentUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayAttachmentDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayAttachment, resourceNsxtPolicyTransitGatewayAttachmentDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayBfdPeerCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayBfdPeer, resourceNsxtPolicyTransitGatewayBfdPeerCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayBfdPeerRead_invalid_parent_path", resourceNsxtPolicyTransitGatewayBfdPeer, resourceNsxtPolicyTransitGatewayBfdPeerRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayBfdPeerUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayBfdPeer, resourceNsxtPolicyTransitGatewayBfdPeerUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayBfdPeerDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayBfdPeer, resourceNsxtPolicyTransitGatewayBfdPeerDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayBgpNeighborCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayBgpNeighbor, resourceNsxtPolicyTransitGatewayBgpNeighborCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayBgpNeighborRead_invalid_parent_path", resourceNsxtPolicyTransitGatewayBgpNeighbor, resourceNsxtPolicyTransitGatewayBgpNeighborRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayBgpNeighborUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayBgpNeighbor, resourceNsxtPolicyTransitGatewayBgpNeighborUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayBgpNeighborDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayBgpNeighbor, resourceNsxtPolicyTransitGatewayBgpNeighborDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayCommunityListCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayCommunityList, resourceNsxtPolicyTransitGatewayCommunityListCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayCommunityListRead_invalid_parent_path", resourceNsxtPolicyTransitGatewayCommunityList, resourceNsxtPolicyTransitGatewayCommunityListRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayCommunityListUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayCommunityList, resourceNsxtPolicyTransitGatewayCommunityListUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayCommunityListDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayCommunityList, resourceNsxtPolicyTransitGatewayCommunityListDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTGWIPSecVpnLocalEndpointCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint, resourceNsxtPolicyTGWIPSecVpnLocalEndpointCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTGWIPSecVpnLocalEndpointUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint, resourceNsxtPolicyTGWIPSecVpnLocalEndpointUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTGWIPSecVpnLocalEndpointDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint, resourceNsxtPolicyTGWIPSecVpnLocalEndpointDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTGWIPSecVpnServicesCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayIpsecVpnService, resourceNsxtPolicyTGWIPSecVpnServicesCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTGWIPSecVpnServicesUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayIpsecVpnService, resourceNsxtPolicyTGWIPSecVpnServicesUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTGWIPSecVpnServicesDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayIpsecVpnService, resourceNsxtPolicyTGWIPSecVpnServicesDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTGWIPSecVpnSessionCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayIPSecVpnSession, resourceNsxtPolicyTGWIPSecVpnSessionCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTGWIPSecVpnSessionUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayIPSecVpnSession, resourceNsxtPolicyTGWIPSecVpnSessionUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTGWIPSecVpnSessionDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayIPSecVpnSession, resourceNsxtPolicyTGWIPSecVpnSessionDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayNatRuleCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayNatRule, resourceNsxtPolicyTransitGatewayNatRuleCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayNatRuleRead_invalid_parent_path", resourceNsxtPolicyTransitGatewayNatRule, resourceNsxtPolicyTransitGatewayNatRuleRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayNatRuleUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayNatRule, resourceNsxtPolicyTransitGatewayNatRuleUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayNatRuleDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayNatRule, resourceNsxtPolicyTransitGatewayNatRuleDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayPrefixListCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayPrefixList, resourceNsxtPolicyTransitGatewayPrefixListCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayPrefixListRead_invalid_parent_path", resourceNsxtPolicyTransitGatewayPrefixList, resourceNsxtPolicyTransitGatewayPrefixListRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayPrefixListUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayPrefixList, resourceNsxtPolicyTransitGatewayPrefixListUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayPrefixListDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayPrefixList, resourceNsxtPolicyTransitGatewayPrefixListDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayRouteMapCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayRouteMap, resourceNsxtPolicyTransitGatewayRouteMapCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayRouteMapRead_invalid_parent_path", resourceNsxtPolicyTransitGatewayRouteMap, resourceNsxtPolicyTransitGatewayRouteMapRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayRouteMapUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayRouteMap, resourceNsxtPolicyTransitGatewayRouteMapUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayRouteMapDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayRouteMap, resourceNsxtPolicyTransitGatewayRouteMapDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayStaticRouteCreate_invalid_parent_path", resourceNsxtPolicyTransitGatewayStaticRoute, resourceNsxtPolicyTransitGatewayStaticRouteCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayStaticRouteRead_invalid_parent_path", resourceNsxtPolicyTransitGatewayStaticRoute, resourceNsxtPolicyTransitGatewayStaticRouteRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayStaticRouteUpdate_invalid_parent_path", resourceNsxtPolicyTransitGatewayStaticRoute, resourceNsxtPolicyTransitGatewayStaticRouteUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyTransitGatewayStaticRouteDelete_invalid_parent_path", resourceNsxtPolicyTransitGatewayStaticRoute, resourceNsxtPolicyTransitGatewayStaticRouteDelete, "parent_path", "invalid", false},
		{"resourceNsxtPolicyVirtualNetworkApplianceRead_missing_id", resourceNsxtPolicyVirtualNetworkAppliance, resourceNsxtPolicyVirtualNetworkApplianceRead, "", "", true},
		{"resourceNsxtPolicyVirtualNetworkApplianceUpdate_missing_id", resourceNsxtPolicyVirtualNetworkAppliance, resourceNsxtPolicyVirtualNetworkApplianceUpdate, "", "", true},
		{"resourceNsxtPolicyVirtualNetworkApplianceDelete_missing_id", resourceNsxtPolicyVirtualNetworkAppliance, resourceNsxtPolicyVirtualNetworkApplianceDelete, "", "", true},
		{"resourceNsxtProxyConfigUpdate_missing_id", resourceNsxtProxyConfig, resourceNsxtProxyConfigUpdate, "", "", true},
		{"resourceNsxtProxyConfigDelete_missing_id", resourceNsxtProxyConfig, resourceNsxtProxyConfigDelete, "", "", true},
		{"resourceNsxtVpcRead_missing_id", resourceNsxtVpc, resourceNsxtVpcRead, "", "", true},
		{"resourceNsxtVpcUpdate_missing_id", resourceNsxtVpc, resourceNsxtVpcUpdate, "", "", true},
		{"resourceNsxtVpcDelete_missing_id", resourceNsxtVpc, resourceNsxtVpcDelete, "", "", true},
		{"resourceNsxtVpcAttachmentCreate_invalid_parent_path", resourceNsxtVpcAttachment, resourceNsxtVpcAttachmentCreate, "parent_path", "invalid", false},
		{"resourceNsxtVpcAttachmentRead_missing_id", resourceNsxtVpcAttachment, resourceNsxtVpcAttachmentRead, "", "", true},
		{"resourceNsxtVpcAttachmentRead_invalid_parent_path", resourceNsxtVpcAttachment, resourceNsxtVpcAttachmentRead, "parent_path", "invalid", false},
		{"resourceNsxtVpcAttachmentUpdate_missing_id", resourceNsxtVpcAttachment, resourceNsxtVpcAttachmentUpdate, "", "", true},
		{"resourceNsxtVpcAttachmentUpdate_invalid_parent_path", resourceNsxtVpcAttachment, resourceNsxtVpcAttachmentUpdate, "parent_path", "invalid", false},
		{"resourceNsxtVpcAttachmentDelete_missing_id", resourceNsxtVpcAttachment, resourceNsxtVpcAttachmentDelete, "", "", true},
		{"resourceNsxtVpcAttachmentDelete_invalid_parent_path", resourceNsxtVpcAttachment, resourceNsxtVpcAttachmentDelete, "parent_path", "invalid", false},
		{"resourceNsxtVpcConnectivityProfileRead_missing_id", resourceNsxtVpcConnectivityProfile, resourceNsxtVpcConnectivityProfileRead, "", "", true},
		{"resourceNsxtVpcConnectivityProfileUpdate_missing_id", resourceNsxtVpcConnectivityProfile, resourceNsxtVpcConnectivityProfileUpdate, "", "", true},
		{"resourceNsxtVpcConnectivityProfileDelete_missing_id", resourceNsxtVpcConnectivityProfile, resourceNsxtVpcConnectivityProfileDelete, "", "", true},
		{"resourceNsxtVpcSubnetDhcpV4StaticBindingConfigCreate_invalid_parent_path", resourceNsxtVpcSubnetDhcpV4StaticBindingConfig, resourceNsxtVpcSubnetDhcpV4StaticBindingConfigCreate, "parent_path", "invalid", false},
		{"resourceNsxtVpcSubnetDhcpV4StaticBindingConfigRead_missing_id", resourceNsxtVpcSubnetDhcpV4StaticBindingConfig, resourceNsxtVpcSubnetDhcpV4StaticBindingConfigRead, "", "", true},
		{"resourceNsxtVpcSubnetDhcpV4StaticBindingConfigRead_invalid_parent_path", resourceNsxtVpcSubnetDhcpV4StaticBindingConfig, resourceNsxtVpcSubnetDhcpV4StaticBindingConfigRead, "parent_path", "invalid", false},
		{"resourceNsxtVpcSubnetDhcpV4StaticBindingConfigUpdate_missing_id", resourceNsxtVpcSubnetDhcpV4StaticBindingConfig, resourceNsxtVpcSubnetDhcpV4StaticBindingConfigUpdate, "", "", true},
		{"resourceNsxtVpcSubnetDhcpV4StaticBindingConfigUpdate_invalid_parent_path", resourceNsxtVpcSubnetDhcpV4StaticBindingConfig, resourceNsxtVpcSubnetDhcpV4StaticBindingConfigUpdate, "parent_path", "invalid", false},
		{"resourceNsxtVpcSubnetDhcpV4StaticBindingConfigDelete_missing_id", resourceNsxtVpcSubnetDhcpV4StaticBindingConfig, resourceNsxtVpcSubnetDhcpV4StaticBindingConfigDelete, "", "", true},
		{"resourceNsxtVpcSubnetDhcpV4StaticBindingConfigDelete_invalid_parent_path", resourceNsxtVpcSubnetDhcpV4StaticBindingConfig, resourceNsxtVpcSubnetDhcpV4StaticBindingConfigDelete, "parent_path", "invalid", false},
		{"resourceNsxtVpcSubnetDhcpV6StaticBindingConfigCreate_invalid_parent_path", resourceNsxtVpcSubnetDhcpV6StaticBindingConfig, resourceNsxtVpcSubnetDhcpV6StaticBindingConfigCreate, "parent_path", "invalid", false},
		{"resourceNsxtVpcSubnetDhcpV6StaticBindingConfigRead_missing_id", resourceNsxtVpcSubnetDhcpV6StaticBindingConfig, resourceNsxtVpcSubnetDhcpV6StaticBindingConfigRead, "", "", true},
		{"resourceNsxtVpcSubnetDhcpV6StaticBindingConfigRead_invalid_parent_path", resourceNsxtVpcSubnetDhcpV6StaticBindingConfig, resourceNsxtVpcSubnetDhcpV6StaticBindingConfigRead, "parent_path", "invalid", false},
		{"resourceNsxtVpcSubnetDhcpV6StaticBindingConfigUpdate_missing_id", resourceNsxtVpcSubnetDhcpV6StaticBindingConfig, resourceNsxtVpcSubnetDhcpV6StaticBindingConfigUpdate, "", "", true},
		{"resourceNsxtVpcSubnetDhcpV6StaticBindingConfigUpdate_invalid_parent_path", resourceNsxtVpcSubnetDhcpV6StaticBindingConfig, resourceNsxtVpcSubnetDhcpV6StaticBindingConfigUpdate, "parent_path", "invalid", false},
		{"resourceNsxtVpcSubnetDhcpV6StaticBindingConfigDelete_missing_id", resourceNsxtVpcSubnetDhcpV6StaticBindingConfig, resourceNsxtVpcSubnetDhcpV6StaticBindingConfigDelete, "", "", true},
		{"resourceNsxtVpcSubnetDhcpV6StaticBindingConfigDelete_invalid_parent_path", resourceNsxtVpcSubnetDhcpV6StaticBindingConfig, resourceNsxtVpcSubnetDhcpV6StaticBindingConfigDelete, "parent_path", "invalid", false},
		{"resourceNsxtVpcEndpointRead_missing_id", resourceNsxtVpcEndpoint, resourceNsxtVpcEndpointRead, "", "", true},
		{"resourceNsxtVpcEndpointUpdate_missing_id", resourceNsxtVpcEndpoint, resourceNsxtVpcEndpointUpdate, "", "", true},
		{"resourceNsxtVpcEndpointDelete_missing_id", resourceNsxtVpcEndpoint, resourceNsxtVpcEndpointDelete, "", "", true},
		{"resourceNsxtVPCGroupRead_missing_id", resourceNsxtVPCGroup, resourceNsxtVPCGroupRead, "", "", true},
		{"resourceNsxtVpcIpAddressAllocationRead_missing_id", resourceNsxtVpcIpAddressAllocation, resourceNsxtVpcIpAddressAllocationRead, "", "", true},
		{"resourceNsxtVpcIpAddressAllocationUpdate_missing_id", resourceNsxtVpcIpAddressAllocation, resourceNsxtVpcIpAddressAllocationUpdate, "", "", true},
		{"resourceNsxtVpcIpAddressAllocationDelete_missing_id", resourceNsxtVpcIpAddressAllocation, resourceNsxtVpcIpAddressAllocationDelete, "", "", true},
		{"resourceNsxtPolicyVpcNatRuleCreate_invalid_parent_path", resourceNsxtPolicyVpcNatRule, resourceNsxtPolicyVpcNatRuleCreate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyVpcNatRuleRead_missing_id", resourceNsxtPolicyVpcNatRule, resourceNsxtPolicyVpcNatRuleRead, "", "", true},
		{"resourceNsxtPolicyVpcNatRuleRead_invalid_parent_path", resourceNsxtPolicyVpcNatRule, resourceNsxtPolicyVpcNatRuleRead, "parent_path", "invalid", false},
		{"resourceNsxtPolicyVpcNatRuleUpdate_missing_id", resourceNsxtPolicyVpcNatRule, resourceNsxtPolicyVpcNatRuleUpdate, "", "", true},
		{"resourceNsxtPolicyVpcNatRuleUpdate_invalid_parent_path", resourceNsxtPolicyVpcNatRule, resourceNsxtPolicyVpcNatRuleUpdate, "parent_path", "invalid", false},
		{"resourceNsxtPolicyVpcNatRuleDelete_missing_id", resourceNsxtPolicyVpcNatRule, resourceNsxtPolicyVpcNatRuleDelete, "", "", true},
		{"resourceNsxtPolicyVpcNatRuleDelete_invalid_parent_path", resourceNsxtPolicyVpcNatRule, resourceNsxtPolicyVpcNatRuleDelete, "parent_path", "invalid", false},
		{"resourceNsxtVpcServiceEndpointRead_missing_id", resourceNsxtVpcServiceEndpoint, resourceNsxtVpcServiceEndpointRead, "", "", true},
		{"resourceNsxtVpcServiceEndpointUpdate_missing_id", resourceNsxtVpcServiceEndpoint, resourceNsxtVpcServiceEndpointUpdate, "", "", true},
		{"resourceNsxtVpcServiceEndpointDelete_missing_id", resourceNsxtVpcServiceEndpoint, resourceNsxtVpcServiceEndpointDelete, "", "", true},
		{"resourceNsxtVpcServiceProfileRead_missing_id", resourceNsxtVpcServiceProfile, resourceNsxtVpcServiceProfileRead, "", "", true},
		{"resourceNsxtVpcServiceProfileUpdate_missing_id", resourceNsxtVpcServiceProfile, resourceNsxtVpcServiceProfileUpdate, "", "", true},
		{"resourceNsxtVpcServiceProfileDelete_missing_id", resourceNsxtVpcServiceProfile, resourceNsxtVpcServiceProfileDelete, "", "", true},
		{"resourceNsxtVpcStaticRoutesRead_missing_id", resourceNsxtVpcStaticRoutes, resourceNsxtVpcStaticRoutesRead, "", "", true},
		{"resourceNsxtVpcStaticRoutesUpdate_missing_id", resourceNsxtVpcStaticRoutes, resourceNsxtVpcStaticRoutesUpdate, "", "", true},
		{"resourceNsxtVpcStaticRoutesDelete_missing_id", resourceNsxtVpcStaticRoutes, resourceNsxtVpcStaticRoutesDelete, "", "", true},
		{"resourceNsxtVpcSubnetRead_missing_id", resourceNsxtVpcSubnet, resourceNsxtVpcSubnetRead, "", "", true},
		{"resourceNsxtVpcSubnetUpdate_missing_id", resourceNsxtVpcSubnet, resourceNsxtVpcSubnetUpdate, "", "", true},
		{"resourceNsxtVpcSubnetDelete_missing_id", resourceNsxtVpcSubnet, resourceNsxtVpcSubnetDelete, "", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.resource()
			raw := map[string]interface{}{}
			for k, v := range map[string]interface{}{
				"gateway_path": "/infra/tier-0s/gateway", "segment_path": "/infra/segments/segment", "pool_path": "/infra/ip-pools/pool", "site_path": "/infra/sites/default",
				"parent_path": "/orgs/default/projects/project/vpcs/vpc", "domain": "default", "domain_id": "default", "group_path": "/infra/domains/default/groups/group",
				"locale_service_id": "default", "gateway_id": "gateway", "service_path": "/infra/tier-0s/gateway/ipsec-vpn-services/vpn", "enforcement_point": "default",
			} {
				if _, ok := res.Schema[k]; ok {
					raw[k] = v
				}
			}
			if strings.Contains(tc.name, "Tier1") {
				raw["gateway_path"] = "/infra/tier-1s/gateway"
			}
			if tc.field != "" {
				raw[tc.field] = tc.value
			}
			d := schema.TestResourceDataRaw(t, res.Schema, raw)
			if !tc.missingID && !strings.Contains(tc.name, "Create_") {
				d.SetId("test-id")
			}
			originalID := d.Id()
			m := newIsolatedMockProviderClient(t)
			m.Host = "https://nsx.invalid"
			m.PolicyHTTPClient = &http.Client{Transport: unitPolicyTransport(func(r *http.Request) (*http.Response, error) {
				t.Fatalf("invalid input reached backend: %s %s", r.Method, r.URL)
				return nil, nil
			})}
			err := tc.operation(d, m)
			require.Error(t, err)
			if tc.missingID {
				require.Contains(t, strings.ToLower(err.Error()), "id")
			} else {
				require.Contains(t, strings.ToLower(err.Error()), "path")
			}
			require.Equal(t, originalID, d.Id(), "validation failures must leave state unchanged")
		})
	}
}
