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

	excludelistapi "github.com/vmware/terraform-provider-nsxt/api/infra/settings/firewall/security"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	excludemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/settings/firewall/security"
)

var (
	exclMemberPath = "/infra/domains/default/groups/test-group"
)

func exclListAPIResponse(members []string) nsxModel.PolicyExcludeList {
	return nsxModel.PolicyExcludeList{
		Members: members,
	}
}

func setupExcludeListMock(t *testing.T, ctrl *gomock.Controller) (*excludemocks.MockExcludeListClient, func()) {
	mockSDK := excludemocks.NewMockExcludeListClient(ctrl)
	mockWrapper := &excludelistapi.PolicyExcludeListClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliExcludeListClient
	cliExcludeListClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *excludelistapi.PolicyExcludeListClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliExcludeListClient = original }
}

func TestMockResourceNsxtPolicyFirewallExcludeListMemberCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupExcludeListMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get().Return(exclListAPIResponse([]string{}), nil),
			mockSDK.EXPECT().Update(gomock.Any()).Return(exclListAPIResponse([]string{exclMemberPath}), nil),
			mockSDK.EXPECT().Get().Return(exclListAPIResponse([]string{exclMemberPath}), nil),
		)

		res := resourceNsxtPolicyFirewallExcludeListMember()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"member": exclMemberPath,
		})

		err := resourceNsxtPolicyFirewallExcludeListMemberCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, exclMemberPath, d.Id())
	})

	t.Run("Create fails when already in list", func(t *testing.T) {
		mockSDK.EXPECT().Get().Return(exclListAPIResponse([]string{exclMemberPath}), nil)

		res := resourceNsxtPolicyFirewallExcludeListMember()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"member": exclMemberPath,
		})

		err := resourceNsxtPolicyFirewallExcludeListMemberCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyFirewallExcludeListMemberRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupExcludeListMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get().Return(exclListAPIResponse([]string{exclMemberPath}), nil)

		res := resourceNsxtPolicyFirewallExcludeListMember()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"member": exclMemberPath,
		})
		d.SetId(exclMemberPath)

		err := resourceNsxtPolicyFirewallExcludeListMemberRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, exclMemberPath, d.Get("member"))
	})

	t.Run("Read returns not found when member absent from list", func(t *testing.T) {
		mockSDK.EXPECT().Get().Return(exclListAPIResponse([]string{}), nil)

		res := resourceNsxtPolicyFirewallExcludeListMember()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"member": exclMemberPath,
		})
		d.SetId(exclMemberPath)

		err := resourceNsxtPolicyFirewallExcludeListMemberRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyFirewallExcludeListMemberDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupExcludeListMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get().Return(exclListAPIResponse([]string{exclMemberPath}), nil),
			mockSDK.EXPECT().Update(gomock.Any()).Return(exclListAPIResponse([]string{}), nil),
		)

		res := resourceNsxtPolicyFirewallExcludeListMember()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"member": exclMemberPath,
		})
		d.SetId(exclMemberPath)

		err := resourceNsxtPolicyFirewallExcludeListMemberDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete is no-op when list not found", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		mockSDK.EXPECT().Get().Return(exclListAPIResponse(nil), notFoundErr)

		res := resourceNsxtPolicyFirewallExcludeListMember()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"member": exclMemberPath,
		})
		d.SetId(exclMemberPath)

		err := resourceNsxtPolicyFirewallExcludeListMemberDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
