//go:build unittest

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

	apinat "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs/nat"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	natmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs/nat"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vpcNatRuleID          = "vpc-nat-rule-id"
	vpcNatRuleDisplayName = "test-vpc-nat-rule"
	vpcNatRuleDescription = "Test VPC NAT Rule"
	vpcNatRuleRevision    = int64(1)
	vpcNatRuleParentPath  = "/orgs/default/projects/project1/vpcs/vpc1/nat/USER"
	vpcNatRuleAction      = "SNAT"
)

func vpcNatRuleAPIResponse() nsxModel.PolicyVpcNatRule {
	return nsxModel.PolicyVpcNatRule{
		Id:          &vpcNatRuleID,
		DisplayName: &vpcNatRuleDisplayName,
		Description: &vpcNatRuleDescription,
		Revision:    &vpcNatRuleRevision,
		Action:      &vpcNatRuleAction,
	}
}

func minimalVpcNatRuleData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vpcNatRuleDisplayName,
		"description":  vpcNatRuleDescription,
		"nsx_id":       vpcNatRuleID,
		"parent_path":  vpcNatRuleParentPath,
		"action":       vpcNatRuleAction,
	}
}

func setupNatRuleMock(t *testing.T, ctrl *gomock.Controller) (*natmocks.MockNatRulesClient, func()) {
	mockSDK := natmocks.NewMockNatRulesClient(ctrl)
	mockWrapper := &apinat.VpcNatRuleClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	original := cliVpcNatRulesClient
	cliVpcNatRulesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apinat.VpcNatRuleClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcNatRulesClient = original }
}

func TestMockResourceNsxtVpcNatRuleCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNatRuleMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), vpcNatRuleID).Return(nsxModel.PolicyVpcNatRule{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), vpcNatRuleID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), vpcNatRuleID).Return(vpcNatRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyVpcNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcNatRuleData())

		err := resourceNsxtPolicyVpcNatRuleCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcNatRuleID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyVpcNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcNatRuleData())

		err := resourceNsxtPolicyVpcNatRuleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcNatRuleRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNatRuleMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), vpcNatRuleID).Return(vpcNatRuleAPIResponse(), nil)

		res := resourceNsxtPolicyVpcNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcNatRuleData())
		d.SetId(vpcNatRuleID)

		err := resourceNsxtPolicyVpcNatRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcNatRuleDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNatRuleMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), vpcNatRuleID).Return(nsxModel.PolicyVpcNatRule{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyVpcNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcNatRuleData())
		d.SetId(vpcNatRuleID)

		err := resourceNsxtPolicyVpcNatRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcNatRuleUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNatRuleMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), vpcNatRuleID, gomock.Any()).Return(vpcNatRuleAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), vpcNatRuleID).Return(vpcNatRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyVpcNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcNatRuleData())
		d.SetId(vpcNatRuleID)

		err := resourceNsxtPolicyVpcNatRuleUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcNatRuleDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNatRuleMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), vpcNatRuleID).Return(nil)

		res := resourceNsxtPolicyVpcNatRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcNatRuleData())
		d.SetId(vpcNatRuleID)

		err := resourceNsxtPolicyVpcNatRuleDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
