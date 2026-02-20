// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/tier_0s/nat/NatRulesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/nat/NatRulesClient.go NatRulesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	t0nat "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/nat"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	natmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/nat"
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
