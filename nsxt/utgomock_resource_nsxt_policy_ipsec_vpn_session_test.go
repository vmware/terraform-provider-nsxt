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
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	tier0ipsecvpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/ipsec_vpn_services"
	tier0localeipsecvpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services/ipsec_vpn_services"
	ipsecvpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/ipsec_vpn_services"
	tier1localeipsecvpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/locale_services/ipsec_vpn_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tier0ipsecmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/ipsec_vpn_services"
	tier0localeipsecmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services/ipsec_vpn_services"
	t1IpsecMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/ipsec_vpn_services"
	tier1localeipsecmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/locale_services/ipsec_vpn_services"
)

var (
	ipsecSessionID          = "ipsec-session-1"
	ipsecSessionServicePath = "/infra/tier-1s/t1-gw-1/ipsec-vpn-services/svc-1"
	ipsecSessionGwID        = "t1-gw-1"
	ipsecSessionSvcID       = "svc-1"
)

func minimalIPSecSessionData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":               "Test IPSec Session",
		"description":                "Test ipsec vpn session",
		"nsx_id":                     ipsecSessionID,
		"service_path":               ipsecSessionServicePath,
		"vpn_type":                   routeBasedIPSecVpnSession,
		"peer_id":                    "10.20.30.40",
		"peer_address":               "10.20.30.40",
		"ip_addresses":               []interface{}{"192.168.10.1"},
		"prefix_length":              24,
		"enabled":                    true,
		"compliance_suite":           nsxModel.IPSecVpnSession_COMPLIANCE_SUITE_NONE,
		"authentication_mode":        nsxModel.IPSecVpnSession_AUTHENTICATION_MODE_PSK,
		"connection_initiation_mode": nsxModel.IPSecVpnSession_CONNECTION_INITIATION_MODE_INITIATOR,
	}
}

func setupIPSecSessionMock(t *testing.T, ctrl *gomock.Controller) (*t1IpsecMocks.MockSessionsClient, func()) {
	mockSDK := t1IpsecMocks.NewMockSessionsClient(ctrl)
	mockWrapper := &ipsecvpnapi.IpsecVpnSessionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier1IpsecVpnSessionsClient
	cliTier1IpsecVpnSessionsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *ipsecvpnapi.IpsecVpnSessionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier1IpsecVpnSessionsClient = original }
}

// setupTier0IPSecSessionMock wires the non-locale-service tier-0 client seam
// (isT0 == true, localeServiceID == "").
func setupTier0IPSecSessionMock(t *testing.T, ctrl *gomock.Controller) (*tier0ipsecmocks.MockSessionsClient, func()) {
	mockSDK := tier0ipsecmocks.NewMockSessionsClient(ctrl)
	mockWrapper := &tier0ipsecvpnapi.IpsecVpnSessionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier0IpsecVpnSessionsClient
	cliTier0IpsecVpnSessionsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0ipsecvpnapi.IpsecVpnSessionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier0IpsecVpnSessionsClient = original }
}

// setupTier0LocaleIPSecSessionMock wires the locale-service tier-0 client seam
// (isT0 == true, localeServiceID != "").
func setupTier0LocaleIPSecSessionMock(t *testing.T, ctrl *gomock.Controller) (*tier0localeipsecmocks.MockSessionsClient, func()) {
	mockSDK := tier0localeipsecmocks.NewMockSessionsClient(ctrl)
	mockWrapper := &tier0localeipsecvpnapi.IpsecVpnSessionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier0LocaleServiceIpsecVpnSessionsClient
	cliTier0LocaleServiceIpsecVpnSessionsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0localeipsecvpnapi.IpsecVpnSessionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier0LocaleServiceIpsecVpnSessionsClient = original }
}

// setupTier1LocaleIPSecSessionMock wires the locale-service tier-1 client seam
// (isT0 == false, localeServiceID != "").
func setupTier1LocaleIPSecSessionMock(t *testing.T, ctrl *gomock.Controller) (*tier1localeipsecmocks.MockSessionsClient, func()) {
	mockSDK := tier1localeipsecmocks.NewMockSessionsClient(ctrl)
	mockWrapper := &tier1localeipsecvpnapi.IpsecVpnSessionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier1LocaleServiceIpsecVpnSessionsClient
	cliTier1LocaleServiceIpsecVpnSessionsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier1localeipsecvpnapi.IpsecVpnSessionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier1LocaleServiceIpsecVpnSessionsClient = original }
}

// ipsecRouteBasedStructValue builds a real *data.StructValue for a RouteBasedIPSecVpnSession,
// mimicking what the SDK Get() call returns on the wire.
func ipsecRouteBasedStructValue(t *testing.T, id, displayName string) *vapiData.StructValue {
	converter := bindings.NewTypeConverter()
	peerID := "10.20.30.40"
	peerAddress := "10.20.30.40"
	authMode := nsxModel.IPSecVpnSession_AUTHENTICATION_MODE_PSK
	complianceSuite := nsxModel.IPSecVpnSession_COMPLIANCE_SUITE_NONE
	connMode := nsxModel.IPSecVpnSession_CONNECTION_INITIATION_MODE_INITIATOR
	enabled := true
	psk := "secret"
	prefixLength := int64(24)

	obj := nsxModel.RouteBasedIPSecVpnSession{
		Id:                       &id,
		DisplayName:              &displayName,
		ConnectionInitiationMode: &connMode,
		ComplianceSuite:          &complianceSuite,
		AuthenticationMode:       &authMode,
		ResourceType:             nsxModel.IPSecVpnSession_RESOURCE_TYPE_ROUTEBASEDIPSECVPNSESSION,
		Enabled:                  &enabled,
		PeerAddress:              &peerAddress,
		PeerId:                   &peerID,
		Psk:                      &psk,
		TunnelInterfaces: []nsxModel.IPSecVpnTunnelInterface{
			{
				DisplayName: &displayName,
				IpSubnets: []nsxModel.TunnelInterfaceIPSubnet{
					{
						IpAddresses:  []string{"192.168.10.1"},
						PrefixLength: &prefixLength,
					},
				},
			},
		},
	}
	dataValue, errs := converter.ConvertToVapi(obj, nsxModel.RouteBasedIPSecVpnSessionBindingType())
	require.Nil(t, errs)
	return dataValue.(*vapiData.StructValue)
}

// ipsecPolicyBasedStructValue builds a real *data.StructValue for a PolicyBasedIPSecVpnSession
// with one rule, mimicking what the SDK Get() call returns on the wire.
func ipsecPolicyBasedStructValue(t *testing.T, id, displayName string) *vapiData.StructValue {
	converter := bindings.NewTypeConverter()
	peerID := "10.20.30.40"
	peerAddress := "10.20.30.40"
	authMode := nsxModel.IPSecVpnSession_AUTHENTICATION_MODE_PSK
	complianceSuite := nsxModel.IPSecVpnSession_COMPLIANCE_SUITE_NONE
	connMode := nsxModel.IPSecVpnSession_CONNECTION_INITIATION_MODE_INITIATOR
	enabled := true
	psk := "secret"
	ruleID := "rule-1"
	srcSubnet := "10.0.0.0/24"
	dstSubnet := "20.0.0.0/24"
	action := nsxModel.IPSecVpnRule_ACTION_PROTECT

	obj := nsxModel.PolicyBasedIPSecVpnSession{
		Id:                       &id,
		DisplayName:              &displayName,
		ConnectionInitiationMode: &connMode,
		ComplianceSuite:          &complianceSuite,
		AuthenticationMode:       &authMode,
		ResourceType:             nsxModel.IPSecVpnSession_RESOURCE_TYPE_POLICYBASEDIPSECVPNSESSION,
		Enabled:                  &enabled,
		PeerAddress:              &peerAddress,
		PeerId:                   &peerID,
		Psk:                      &psk,
		Rules: []nsxModel.IPSecVpnRule{
			{
				Id:           &ruleID,
				UniqueId:     &ruleID,
				Sources:      []nsxModel.IPSecVpnSubnet{{Subnet: &srcSubnet}},
				Destinations: []nsxModel.IPSecVpnSubnet{{Subnet: &dstSubnet}},
				Action:       &action,
			},
		},
	}
	dataValue, errs := converter.ConvertToVapi(obj, nsxModel.PolicyBasedIPSecVpnSessionBindingType())
	require.Nil(t, errs)
	return dataValue.(*vapiData.StructValue)
}

func TestMockResourceNsxtPolicyIPSecVpnSessionRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSessionMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnSessionDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSessionMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())

		err := resourceNsxtPolicyIPSecVpnSessionDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnSessionReadSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSessionMock(t, ctrl)
	defer restore()

	t.Run("Read success route-based sets fields", func(t *testing.T) {
		sv := ipsecRouteBasedStructValue(t, ipsecSessionID, "Test IPSec Session")
		mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(sv, nil)

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "Test IPSec Session", d.Get("display_name"))
		assert.Equal(t, ipsecSessionID, d.Get("nsx_id"))
		assert.Equal(t, routeBasedIPSecVpnSession, d.Get("vpn_type"))
		assert.Equal(t, "10.20.30.40", d.Get("peer_id"))
		assert.Equal(t, "10.20.30.40", d.Get("peer_address"))
		assert.Equal(t, 24, d.Get("prefix_length"))
		assert.True(t, d.Get("enabled").(bool))

		ips := d.Get("ip_addresses").([]interface{})
		require.Len(t, ips, 1)
		assert.Equal(t, "192.168.10.1", ips[0])
	})

	t.Run("Read success policy-based sets rules", func(t *testing.T) {
		sv := ipsecPolicyBasedStructValue(t, ipsecSessionID, "Test IPSec Session")
		mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(sv, nil)

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, policyBasedIPSecVpnSession, d.Get("vpn_type"))

		rules := d.Get("rule").([]interface{})
		require.Len(t, rules, 1)
		rule := rules[0].(map[string]interface{})
		assert.Equal(t, "rule-1", rule["nsx_id"])
		assert.Equal(t, nsxModel.IPSecVpnRule_ACTION_PROTECT, rule["action"])
	})
}

func TestMockResourceNsxtPolicyIPSecVpnSessionCreateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSessionMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		sv := ipsecRouteBasedStructValue(t, ipsecSessionID, "Test IPSec Session")

		gomock.InOrder(
			mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(nil, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(sv, nil),
		)

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())

		err := resourceNsxtPolicyIPSecVpnSessionCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecSessionID, d.Id())
		assert.Equal(t, ipsecSessionID, d.Get("nsx_id"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		sv := ipsecRouteBasedStructValue(t, ipsecSessionID, "Test IPSec Session")
		mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(sv, nil)

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())

		err := resourceNsxtPolicyIPSecVpnSessionCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when Patch API errors", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(nil, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())

		err := resourceNsxtPolicyIPSecVpnSessionCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnSessionUpdateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSessionMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		sv := ipsecRouteBasedStructValue(t, ipsecSessionID, "Updated Session")

		gomock.InOrder(
			mockSDK.EXPECT().Patch(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(sv, nil),
		)

		res := resourceNsxtPolicyIPSecVpnSession()
		data := minimalIPSecSessionData()
		data["display_name"] = "Updated Session"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "Updated Session", d.Get("display_name"))
	})

	t.Run("Update fails when Patch API errors", func(t *testing.T) {
		mockSDK.EXPECT().Patch(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnSessionTier0NonLocale(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTier0IPSecSessionMock(t, ctrl)
	defer restore()

	t.Run("Read success on tier-0 non-locale-service service path", func(t *testing.T) {
		servicePath := "/infra/tier-0s/t0-gw-1/ipsec-vpn-services/svc-1"
		sv := ipsecRouteBasedStructValue(t, ipsecSessionID, "T0 Session")
		mockSDK.EXPECT().Get("t0-gw-1", "svc-1", ipsecSessionID).Return(sv, nil)

		res := resourceNsxtPolicyIPSecVpnSession()
		data := minimalIPSecSessionData()
		data["service_path"] = servicePath
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "T0 Session", d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyIPSecVpnSessionTier0Locale(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTier0LocaleIPSecSessionMock(t, ctrl)
	defer restore()

	t.Run("Read success on tier-0 locale-service service path", func(t *testing.T) {
		servicePath := "/infra/tier-0s/t0-gw-1/locale-services/default/ipsec-vpn-services/svc-1"
		sv := ipsecRouteBasedStructValue(t, ipsecSessionID, "T0 Locale Session")
		mockSDK.EXPECT().Get("t0-gw-1", "default", "svc-1", ipsecSessionID).Return(sv, nil)

		res := resourceNsxtPolicyIPSecVpnSession()
		data := minimalIPSecSessionData()
		data["service_path"] = servicePath
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "T0 Locale Session", d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyIPSecVpnSessionTier1Locale(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTier1LocaleIPSecSessionMock(t, ctrl)
	defer restore()

	t.Run("Read success on tier-1 locale-service service path", func(t *testing.T) {
		servicePath := "/infra/tier-1s/t1-gw-1/locale-services/default/ipsec-vpn-services/svc-1"
		sv := ipsecRouteBasedStructValue(t, ipsecSessionID, "T1 Locale Session")
		mockSDK.EXPECT().Get("t1-gw-1", "default", "svc-1", ipsecSessionID).Return(sv, nil)

		res := resourceNsxtPolicyIPSecVpnSession()
		data := minimalIPSecSessionData()
		data["service_path"] = servicePath
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "T1 Locale Session", d.Get("display_name"))
	})

	t.Run("Read fails with project context on tier-1 locale-service path", func(t *testing.T) {
		servicePath := "/orgs/default/projects/proj-1/infra/tier-1s/t1-gw-1/locale-services/default/ipsec-vpn-services/svc-1"

		res := resourceNsxtPolicyIPSecVpnSession()
		data := minimalIPSecSessionData()
		data["service_path"] = servicePath
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "project context is not supported")
	})
}
