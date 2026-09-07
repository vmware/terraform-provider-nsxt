//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/tier_0s/nat/NatRulesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/nat/NatRulesClient.go NatRulesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	t0nat "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/nat"
	t1nat "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/nat"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	natmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/nat"
	t1natmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/nat"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	natGatewayPath = "/infra/tier-0s/gw1"
	natGwID        = "gw1"
	natRuleID      = "nat-rule-1"
	natDisplayName = "nat-rule-fooname"
	natDescription = "nat rule mock"
	natAction      = model.PolicyNatRule_ACTION_SNAT
	natType        = model.PolicyNat_NAT_TYPE_USER
)

func TestMockResourceNsxtPolicyNATRuleCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNatSDK := natmocks.NewMockNatRulesClient(ctrl)
	mockWrapper := &t0nat.PolicyNatRuleClientContext{
		Client:     mockNatSDK,
		ClientType: utl.Local,
	}

	originalCli := cliTier0NatRulesClient
	defer func() { cliTier0NatRulesClient = originalCli }()
	cliTier0NatRulesClient = func(sessionContext utl.SessionContext, connector client.Connector) *t0nat.PolicyNatRuleClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockNatSDK.EXPECT().Patch(natGwID, natType, gomock.Any(), gomock.Any()).Return(nil)
		mockNatSDK.EXPECT().Get(natGwID, natType, gomock.Any()).Return(model.PolicyNatRule{
			Id:          &natRuleID,
			DisplayName: &natDisplayName,
			Description: &natDescription,
			Action:      &natAction,
		}, nil)

		res := resourceNsxtPolicyNATRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNATRuleData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNATRuleCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})
}

func TestMockResourceNsxtPolicyNATRuleRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNatSDK := natmocks.NewMockNatRulesClient(ctrl)
	mockWrapper := &t0nat.PolicyNatRuleClientContext{
		Client:     mockNatSDK,
		ClientType: utl.Local,
	}

	originalCli := cliTier0NatRulesClient
	defer func() { cliTier0NatRulesClient = originalCli }()
	cliTier0NatRulesClient = func(sessionContext utl.SessionContext, connector client.Connector) *t0nat.PolicyNatRuleClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockNatSDK.EXPECT().Get(natGwID, natType, natRuleID).Return(model.PolicyNatRule{
			Id:          &natRuleID,
			DisplayName: &natDisplayName,
			Description: &natDescription,
			Action:      &natAction,
		}, nil)

		res := resourceNsxtPolicyNATRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNATRuleData())
		d.SetId(natRuleID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNATRuleRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, natDisplayName, d.Get("display_name"))
		assert.Equal(t, natDescription, d.Get("description"))
		assert.Equal(t, natRuleID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyNATRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNATRuleData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNATRuleRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining NAT Rule ID")
	})
}

func TestMockResourceNsxtPolicyNATRuleUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNatSDK := natmocks.NewMockNatRulesClient(ctrl)
	mockWrapper := &t0nat.PolicyNatRuleClientContext{
		Client:     mockNatSDK,
		ClientType: utl.Local,
	}

	originalCli := cliTier0NatRulesClient
	defer func() { cliTier0NatRulesClient = originalCli }()
	cliTier0NatRulesClient = func(sessionContext utl.SessionContext, connector client.Connector) *t0nat.PolicyNatRuleClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockNatSDK.EXPECT().Patch(natGwID, natType, natRuleID, gomock.Any()).Return(nil)
		mockNatSDK.EXPECT().Get(natGwID, natType, natRuleID).Return(model.PolicyNatRule{
			Id:          &natRuleID,
			DisplayName: &natDisplayName,
			Description: &natDescription,
			Action:      &natAction,
		}, nil)

		res := resourceNsxtPolicyNATRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNATRuleData())
		d.SetId(natRuleID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNATRuleUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyNATRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNATRuleData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNATRuleUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining NAT Rule ID")
	})
}

func TestMockResourceNsxtPolicyNATRuleDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNatSDK := natmocks.NewMockNatRulesClient(ctrl)
	mockWrapper := &t0nat.PolicyNatRuleClientContext{
		Client:     mockNatSDK,
		ClientType: utl.Local,
	}

	originalCli := cliTier0NatRulesClient
	defer func() { cliTier0NatRulesClient = originalCli }()
	cliTier0NatRulesClient = func(sessionContext utl.SessionContext, connector client.Connector) *t0nat.PolicyNatRuleClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockNatSDK.EXPECT().Delete(natGwID, natType, natRuleID).Return(nil)

		res := resourceNsxtPolicyNATRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNATRuleData())
		d.SetId(natRuleID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNATRuleDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyNATRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNATRuleData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNATRuleDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining NAT Rule ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockNatSDK.EXPECT().Delete(natGwID, natType, natRuleID).Return(errors.New("API error"))

		res := resourceNsxtPolicyNATRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNATRuleData())
		d.SetId(natRuleID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNATRuleDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestUnitNsxt_getNatTypeByAction(t *testing.T) {
	assert.Equal(t, model.PolicyNat_NAT_TYPE_NAT64, getNatTypeByAction("", model.PolicyNatRule_ACTION_NAT64))
	assert.Equal(t, model.PolicyNat_NAT_TYPE_USER, getNatTypeByAction("", model.PolicyNatRule_ACTION_SNAT))
	assert.Equal(t, model.PolicyNat_NAT_TYPE_INTERNAL, getNatTypeByAction(model.PolicyNat_NAT_TYPE_INTERNAL, model.PolicyNatRule_ACTION_SNAT))
}

func TestUnitNsxt_translatedNetworksNeeded(t *testing.T) {
	assert.False(t, translatedNetworksNeeded(model.PolicyNatRule_ACTION_NO_SNAT))
	assert.False(t, translatedNetworksNeeded(model.PolicyNatRule_ACTION_NO_DNAT))
	assert.True(t, translatedNetworksNeeded(model.PolicyNatRule_ACTION_SNAT))
}

func TestUnitNsxt_getTranslatedNetworks(t *testing.T) {
	t.Run("errors when required and missing", func(t *testing.T) {
		action := model.PolicyNatRule_ACTION_SNAT
		_, err := getTranslatedNetworks(model.PolicyNatRule{Action: &action})
		require.Error(t, err)
	})

	t.Run("nil is fine when not needed", func(t *testing.T) {
		action := model.PolicyNatRule_ACTION_NO_SNAT
		val, err := getTranslatedNetworks(model.PolicyNatRule{Action: &action})
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("passes through when set", func(t *testing.T) {
		action := model.PolicyNatRule_ACTION_SNAT
		nets := "10.0.0.0/24"
		val, err := getTranslatedNetworks(model.PolicyNatRule{Action: &action, TranslatedNetwork: &nets})
		require.NoError(t, err)
		assert.Equal(t, &nets, val)
	})
}

func TestUnitNsxt_policyBasedVpnModeNeeded(t *testing.T) {
	assert.True(t, policyBasedVpnModeNeeded(model.PolicyNatRule_ACTION_DNAT))
	assert.True(t, policyBasedVpnModeNeeded(model.PolicyNatRule_ACTION_NO_DNAT))
	assert.False(t, policyBasedVpnModeNeeded(model.PolicyNatRule_ACTION_SNAT))
}

func TestUnitNsxt_getPolicyBasedVpnMode(t *testing.T) {
	t.Run("errors when set on unsupported action", func(t *testing.T) {
		action := model.PolicyNatRule_ACTION_SNAT
		match := model.PolicyNatRule_POLICY_BASED_VPN_MODE_BYPASS
		_, err := getPolicyBasedVpnMode(model.PolicyNatRule{Action: &action, PolicyBasedVpnMode: &match})
		require.Error(t, err)
	})

	t.Run("fine when unset", func(t *testing.T) {
		action := model.PolicyNatRule_ACTION_SNAT
		val, err := getPolicyBasedVpnMode(model.PolicyNatRule{Action: &action})
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("fine when set on a supported action", func(t *testing.T) {
		action := model.PolicyNatRule_ACTION_DNAT
		match := model.PolicyNatRule_POLICY_BASED_VPN_MODE_MATCH
		val, err := getPolicyBasedVpnMode(model.PolicyNatRule{Action: &action, PolicyBasedVpnMode: &match})
		require.NoError(t, err)
		assert.Equal(t, &match, val)
	})
}

func TestUnitNsxt_validateNatTypeAction(t *testing.T) {
	t.Run("NAT64 action requires NAT64 type", func(t *testing.T) {
		err := validateNatTypeAction(model.PolicyNatRule_ACTION_NAT64, model.PolicyNat_NAT_TYPE_USER)
		require.Error(t, err)
	})

	t.Run("NAT64 type requires NAT64 action", func(t *testing.T) {
		err := validateNatTypeAction(model.PolicyNatRule_ACTION_SNAT, model.PolicyNat_NAT_TYPE_NAT64)
		require.Error(t, err)
	})

	t.Run("matching NAT64 action and type is fine", func(t *testing.T) {
		require.NoError(t, validateNatTypeAction(model.PolicyNatRule_ACTION_NAT64, model.PolicyNat_NAT_TYPE_NAT64))
	})

	t.Run("non-NAT64 combination is fine", func(t *testing.T) {
		require.NoError(t, validateNatTypeAction(model.PolicyNatRule_ACTION_SNAT, model.PolicyNat_NAT_TYPE_USER))
	})
}

func TestMockResourceNsxtPolicyNATRuleInvalidGatewayPath(t *testing.T) {
	res := resourceNsxtPolicyNATRule()

	for _, fn := range []func(*schema.ResourceData, interface{}) error{
		resourceNsxtPolicyNATRuleCreate,
		resourceNsxtPolicyNATRuleRead,
		resourceNsxtPolicyNATRuleUpdate,
		resourceNsxtPolicyNATRuleDelete,
	} {
		data := minimalNATRuleData()
		data["gateway_path"] = "invalid"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(natRuleID)

		err := fn(d, newGoMockProviderClient())
		require.Error(t, err)
	}
}

func TestMockResourceNsxtPolicyNATRuleTier1Delete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockT1NatSDK := t1natmocks.NewMockNatRulesClient(ctrl)
	mockT1Wrapper := &t1nat.PolicyNatRuleClientContext{
		Client:     mockT1NatSDK,
		ClientType: utl.Local,
	}
	originalT1Cli := cliTier1NatRulesClient
	defer func() { cliTier1NatRulesClient = originalT1Cli }()
	cliTier1NatRulesClient = func(sessionContext utl.SessionContext, connector client.Connector) *t1nat.PolicyNatRuleClientContext {
		return mockT1Wrapper
	}

	t.Run("Delete on a tier-1 gateway path", func(t *testing.T) {
		mockT1NatSDK.EXPECT().Delete("gw1", natType, natRuleID).Return(nil)

		res := resourceNsxtPolicyNATRule()
		data := minimalNATRuleData()
		data["gateway_path"] = "/infra/tier-1s/gw1"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(natRuleID)

		err := resourceNsxtPolicyNATRuleDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyNATRuleCreateAlreadyExists(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNatSDK := natmocks.NewMockNatRulesClient(ctrl)
	mockWrapper := &t0nat.PolicyNatRuleClientContext{
		Client:     mockNatSDK,
		ClientType: utl.Local,
	}
	originalCli := cliTier0NatRulesClient
	defer func() { cliTier0NatRulesClient = originalCli }()
	cliTier0NatRulesClient = func(sessionContext utl.SessionContext, connector client.Connector) *t0nat.PolicyNatRuleClientContext {
		return mockWrapper
	}

	t.Run("Create fails when nsx_id already exists", func(t *testing.T) {
		mockNatSDK.EXPECT().Get(natGwID, natType, "existing-id").Return(model.PolicyNatRule{Id: &natRuleID}, nil)

		res := resourceNsxtPolicyNATRule()
		data := minimalNATRuleData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyNATRuleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func minimalNATRuleData() map[string]interface{} {
	// SNAT requires translated_networks; use a valid CIDR for the test
	return map[string]interface{}{
		"gateway_path":          natGatewayPath,
		"display_name":          natDisplayName,
		"description":           natDescription,
		"action":                natAction,
		"type":                  natType,
		"enabled":               true,
		"firewall_match":        model.PolicyNatRule_FIREWALL_MATCH_BYPASS,
		"logging":               false,
		"rule_priority":         100,
		"service":               "",
		"destination_networks":  []interface{}{},
		"source_networks":       []interface{}{},
		"translated_networks":   []interface{}{"192.168.1.0/24"},
		"translated_ports":      "",
		"scope":                 []interface{}{},
		"policy_based_vpn_mode": "",
	}
}
