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

	dnssvcsapi "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/dns_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	dnssvcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/dns_services"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dnsRuleID          = "dns-rule-id"
	dnsRuleDisplayName = "test-dns-rule"
	dnsRuleDescription = "Test Project DNS Rule"
	dnsRuleRevision    = int64(1)
	dnsRuleParentPath  = "/orgs/default/projects/project1/dns-services/dns-svc-1"
)

func dnsRuleAPIResponse() nsxModel.ProjectDnsRule {
	actionType := nsxModel.ProjectDnsRule_ACTION_TYPE_FORWARD
	upstream := "8.8.8.8"
	return nsxModel.ProjectDnsRule{
		Id:              &dnsRuleID,
		DisplayName:     &dnsRuleDisplayName,
		Description:     &dnsRuleDescription,
		Revision:        &dnsRuleRevision,
		ActionType:      &actionType,
		DomainPatterns:  []string{"*.example.com"},
		UpstreamServers: []string{upstream},
	}
}

func minimalDnsRuleData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":     dnsRuleDisplayName,
		"description":      dnsRuleDescription,
		"nsx_id":           dnsRuleID,
		"parent_path":      dnsRuleParentPath,
		"action_type":      "FORWARD",
		"domain_patterns":  []interface{}{"*.example.com"},
		"upstream_servers": []interface{}{"8.8.8.8"},
	}
}

func setupDnsRuleMock(t *testing.T, ctrl *gomock.Controller) (*dnssvcsmocks.MockRulesClient, func()) {
	mockSDK := dnssvcsmocks.NewMockRulesClient(ctrl)
	mockWrapper := &dnssvcsapi.ProjectDnsRuleClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "project1",
	}

	original := cliProjectDnsRulesClient
	cliProjectDnsRulesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *dnssvcsapi.ProjectDnsRuleClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliProjectDnsRulesClient = original }
}

func TestMockResourceNsxtPolicyDnsRuleCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRuleMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), dnsRuleID).Return(nsxModel.ProjectDnsRule{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), dnsRuleID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), dnsRuleID).Return(dnsRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRuleData())

		err := resourceNsxtPolicyDnsRuleCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsRuleID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyDnsRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRuleData())

		err := resourceNsxtPolicyDnsRuleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})

	t.Run("Create FORWARD rule without upstream_servers fails", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRuleMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), dnsRuleID).Return(nsxModel.ProjectDnsRule{}, notFoundErr)

		res := resourceNsxtPolicyDnsRule()
		data := minimalDnsRuleData()
		data["upstream_servers"] = []interface{}{}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyDnsRuleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "upstream_servers")
	})
}

func TestMockResourceNsxtPolicyDnsRuleRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRuleMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), dnsRuleID).Return(dnsRuleAPIResponse(), nil)

		res := resourceNsxtPolicyDnsRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRuleData())
		d.SetId(dnsRuleID)

		err := resourceNsxtPolicyDnsRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsRuleDisplayName, d.Get("display_name"))
		assert.Equal(t, "FORWARD", d.Get("action_type"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRuleMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), dnsRuleID).Return(nsxModel.ProjectDnsRule{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyDnsRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRuleData())
		d.SetId(dnsRuleID)

		err := resourceNsxtPolicyDnsRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyDnsRuleUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRuleMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), dnsRuleID, gomock.Any()).Return(dnsRuleAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), dnsRuleID).Return(dnsRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRuleData())
		d.SetId(dnsRuleID)

		err := resourceNsxtPolicyDnsRuleUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyDnsRuleDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRuleMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), dnsRuleID).Return(nil)

		res := resourceNsxtPolicyDnsRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRuleData())
		d.SetId(dnsRuleID)

		err := resourceNsxtPolicyDnsRuleDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
