// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apinat "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways/nat"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tgwnatmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways/nat"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	tgwNatRuleID          = "tgw-nat-rule-id"
	tgwNatRuleDisplayName = "test-tgw-nat-rule"
	tgwNatRuleDescription = "Test TGW NAT Rule"
	tgwNatRuleRevision    = int64(1)
	tgwNatRuleParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1/nat/USER"
	tgwNatRuleOrgID       = "default"
	tgwNatRuleProjectID   = "project1"
	tgwNatRuleTGWID       = "tgw1"
	tgwNatRuleNatID       = "USER"
	tgwNatRuleAction      = "SNAT"
)

func tgwNatRuleAPIResponse() nsxModel.TransitGatewayNatRule {
	return nsxModel.TransitGatewayNatRule{
		Id:          &tgwNatRuleID,
		DisplayName: &tgwNatRuleDisplayName,
		Description: &tgwNatRuleDescription,
		Revision:    &tgwNatRuleRevision,
		Action:      &tgwNatRuleAction,
	}
}

func minimalTGWNatRuleData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":   tgwNatRuleDisplayName,
		"description":    tgwNatRuleDescription,
		"nsx_id":         tgwNatRuleID,
		"parent_path":    tgwNatRuleParentPath,
		"action":         tgwNatRuleAction,
		"firewall_match": nsxModel.PolicyVpcNatRule_FIREWALL_MATCH_MATCH_INTERNAL_ADDRESS,
	}
}

func setupTGWNatRuleMock(t *testing.T, ctrl *gomock.Controller) (*tgwnatmocks.MockNatRulesClient, func()) {
	mockSDK := tgwnatmocks.NewMockNatRulesClient(ctrl)
	mockWrapper := &apinat.TransitGatewayNatRuleClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwNatRuleProjectID,
	}

	original := cliTransitGatewayNatRulesClient
	cliTransitGatewayNatRulesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apinat.TransitGatewayNatRuleClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTransitGatewayNatRulesClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayNatRuleCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWNatRuleMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwNatRuleOrgID, tgwNatRuleProjectID, tgwNatRuleTGWID, tgwNatRuleNatID, tgwNatRuleID).Return(nsxModel.TransitGatewayNatRule{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwNatRuleOrgID, tgwNatRuleProjectID, tgwNatRuleTGWID, tgwNatRuleNatID, tgwNatRuleID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwNatRuleOrgID, tgwNatRuleProjectID, tgwNatRuleTGWID, tgwNatRuleNatID, tgwNatRuleID).Return(tgwNatRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWNatRuleData())

		err := resourceNsxtPolicyTransitGatewayNatRuleCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwNatRuleID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyTransitGatewayNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWNatRuleData())

		err := resourceNsxtPolicyTransitGatewayNatRuleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.0")
	})

	t.Run("Create fails with invalid parent_path", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyTransitGatewayNatRule()
		data := minimalTGWNatRuleData()
		data["parent_path"] = "/orgs/default"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayNatRuleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayNatRuleRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWNatRuleMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwNatRuleOrgID, tgwNatRuleProjectID, tgwNatRuleTGWID, tgwNatRuleNatID, tgwNatRuleID).Return(tgwNatRuleAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGatewayNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWNatRuleData())
		d.SetId(tgwNatRuleID)

		err := resourceNsxtPolicyTransitGatewayNatRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwNatRuleDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwNatRuleOrgID, tgwNatRuleProjectID, tgwNatRuleTGWID, tgwNatRuleNatID, tgwNatRuleID).Return(nsxModel.TransitGatewayNatRule{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWNatRuleData())
		d.SetId(tgwNatRuleID)

		err := resourceNsxtPolicyTransitGatewayNatRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWNatRuleData())

		err := resourceNsxtPolicyTransitGatewayNatRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayNatRuleUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWNatRuleMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(tgwNatRuleOrgID, tgwNatRuleProjectID, tgwNatRuleTGWID, tgwNatRuleNatID, tgwNatRuleID, gomock.Any()).Return(tgwNatRuleAPIResponse(), nil),
			mockSDK.EXPECT().Get(tgwNatRuleOrgID, tgwNatRuleProjectID, tgwNatRuleTGWID, tgwNatRuleNatID, tgwNatRuleID).Return(tgwNatRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWNatRuleData())
		d.SetId(tgwNatRuleID)

		err := resourceNsxtPolicyTransitGatewayNatRuleUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWNatRuleData())

		err := resourceNsxtPolicyTransitGatewayNatRuleUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayNatRuleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWNatRuleMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwNatRuleOrgID, tgwNatRuleProjectID, tgwNatRuleTGWID, tgwNatRuleNatID, tgwNatRuleID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWNatRuleData())
		d.SetId(tgwNatRuleID)

		err := resourceNsxtPolicyTransitGatewayNatRuleDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWNatRuleData())

		err := resourceNsxtPolicyTransitGatewayNatRuleDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
