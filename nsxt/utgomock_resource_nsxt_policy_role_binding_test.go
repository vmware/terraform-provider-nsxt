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

	aaaapi "github.com/vmware/terraform-provider-nsxt/api/aaa"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	aaamocks "github.com/vmware/terraform-provider-nsxt/mocks/aaa"
)

var (
	rbID          = "rb-001"
	rbDisplayName = "Test Role Binding"
	rbDescription = "Test role binding description"
	rbRevision    = int64(1)
	rbName        = "testuser@example.com"
	rbType        = nsxModel.RoleBinding_TYPE_REMOTE_USER
)

func rbAPIResponse() nsxModel.RoleBinding {
	return nsxModel.RoleBinding{
		Id:            &rbID,
		DisplayName:   &rbDisplayName,
		Description:   &rbDescription,
		Revision:      &rbRevision,
		Name:          &rbName,
		Type_:         &rbType,
		RolesForPaths: []nsxModel.RolesForPath{},
	}
}

func minimalRbData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":         rbDisplayName,
		"description":          rbDescription,
		"name":                 rbName,
		"type":                 rbType,
		"overwrite_local_user": false,
	}
}

func setupRbMock(t *testing.T, ctrl *gomock.Controller) (*aaamocks.MockRoleBindingsClient, func()) {
	mockSDK := aaamocks.NewMockRoleBindingsClient(ctrl)
	mockWrapper := &aaaapi.RoleBindingClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliRoleBindingsClient
	cliRoleBindingsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *aaaapi.RoleBindingClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliRoleBindingsClient = original }
}

func TestMockResourceNsxtPolicyRoleBindingCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRbMock(t, ctrl)
	defer restore()

	t.Run("Create success for remote user", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Create(gomock.Any()).Return(rbAPIResponse(), nil),
			mockSDK.EXPECT().Get(rbID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(rbAPIResponse(), nil),
		)

		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())

		err := resourceNsxtPolicyUserManagementRoleBindingCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, rbID, d.Id())
	})

	t.Run("Create fails for local user without overwrite flag", func(t *testing.T) {
		localType := nsxModel.RoleBinding_TYPE_LOCAL_USER
		localData := map[string]interface{}{
			"display_name":         rbDisplayName,
			"description":          rbDescription,
			"name":                 rbName,
			"type":                 localType,
			"overwrite_local_user": false,
		}

		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, localData)

		err := resourceNsxtPolicyUserManagementRoleBindingCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "overwrite_local_user")
	})
}

func TestMockResourceNsxtPolicyRoleBindingRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRbMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(rbID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(rbAPIResponse(), nil)

		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		d.SetId(rbID)

		err := resourceNsxtPolicyUserManagementRoleBindingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, rbDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(rbID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.RoleBinding{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		d.SetId(rbID)

		err := resourceNsxtPolicyUserManagementRoleBindingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())

		err := resourceNsxtPolicyUserManagementRoleBindingRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyRoleBindingUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRbMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(rbID, gomock.Any()).Return(rbAPIResponse(), nil),
			mockSDK.EXPECT().Get(rbID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(rbAPIResponse(), nil),
		)

		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		d.SetId(rbID)

		err := resourceNsxtPolicyUserManagementRoleBindingUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())

		err := resourceNsxtPolicyUserManagementRoleBindingUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyRoleBindingDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRbMock(t, ctrl)
	defer restore()

	t.Run("Delete success for remote user", func(t *testing.T) {
		mockSDK.EXPECT().Delete(rbID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(nil)

		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		d.SetId(rbID)

		err := resourceNsxtPolicyUserManagementRoleBindingDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails for local user without overwrite flag", func(t *testing.T) {
		localType := nsxModel.RoleBinding_TYPE_LOCAL_USER
		localData := map[string]interface{}{
			"display_name":         rbDisplayName,
			"description":          rbDescription,
			"name":                 rbName,
			"type":                 localType,
			"overwrite_local_user": false,
		}

		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, localData)
		d.SetId(rbID)

		err := resourceNsxtPolicyUserManagementRoleBindingDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can not be deleted")
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())

		err := resourceNsxtPolicyUserManagementRoleBindingDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
