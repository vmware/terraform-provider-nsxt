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

	aaaapi "github.com/vmware/terraform-provider-nsxt/api/aaa"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	aaamocks "github.com/vmware/terraform-provider-nsxt/mocks/aaa"
)

var (
	roleID          = "my-custom-role"
	roleDisplayName = "My Custom Role"
	roleDescription = "A custom role for testing"
	roleRevision    = int64(1)
	roleRoleName    = "my-custom-role"
)

func roleAPIResponse() nsxModel.RoleWithFeatures {
	return nsxModel.RoleWithFeatures{
		Id:          &roleID,
		DisplayName: &roleDisplayName,
		Description: &roleDescription,
		Revision:    &roleRevision,
		Role:        &roleRoleName,
		Features:    []nsxModel.FeaturePermission{},
	}
}

func minimalRoleData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": roleDisplayName,
		"description":  roleDescription,
		"role":         roleRoleName,
	}
}

func setupRoleMock(t *testing.T, ctrl *gomock.Controller) (*aaamocks.MockRolesClient, func()) {
	mockSDK := aaamocks.NewMockRolesClient(ctrl)
	mockWrapper := &aaaapi.RoleWithFeaturesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliRolesClient
	cliRolesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *aaaapi.RoleWithFeaturesClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliRolesClient = original }
}

func TestMockResourceNsxtPolicyRoleCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRoleMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		emptyResult := nsxModel.RecommendedFeaturePermissionListResult{Results: []nsxModel.RecommendedFeaturePermission{}}
		gomock.InOrder(
			mockSDK.EXPECT().Get(roleRoleName).Return(nsxModel.RoleWithFeatures{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Validate(gomock.Any()).Return(emptyResult, nil),
			mockSDK.EXPECT().Update(roleRoleName, gomock.Any()).Return(nsxModel.RoleWithFeatures{}, nil),
			mockSDK.EXPECT().Get(roleRoleName).Return(roleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyUserManagementRole()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleData())

		err := resourceNsxtPolicyUserManagementRoleCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, roleRoleName, d.Id())
		assert.Equal(t, roleDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(roleRoleName).Return(roleAPIResponse(), nil)

		res := resourceNsxtPolicyUserManagementRole()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleData())

		err := resourceNsxtPolicyUserManagementRoleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyRoleRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRoleMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(roleID).Return(roleAPIResponse(), nil)

		res := resourceNsxtPolicyUserManagementRole()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleData())
		d.SetId(roleID)

		err := resourceNsxtPolicyUserManagementRoleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, roleDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(roleID).Return(nsxModel.RoleWithFeatures{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyUserManagementRole()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleData())
		d.SetId(roleID)

		err := resourceNsxtPolicyUserManagementRoleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyUserManagementRole()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleData())

		err := resourceNsxtPolicyUserManagementRoleRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyRoleUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRoleMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		emptyResult := nsxModel.RecommendedFeaturePermissionListResult{Results: []nsxModel.RecommendedFeaturePermission{}}
		gomock.InOrder(
			mockSDK.EXPECT().Validate(gomock.Any()).Return(emptyResult, nil),
			mockSDK.EXPECT().Update(roleID, gomock.Any()).Return(nsxModel.RoleWithFeatures{}, nil),
			mockSDK.EXPECT().Get(roleID).Return(roleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyUserManagementRole()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleData())
		d.SetId(roleID)

		err := resourceNsxtPolicyUserManagementRoleUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyUserManagementRole()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleData())

		err := resourceNsxtPolicyUserManagementRoleUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyRoleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRoleMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(roleID).Return(nil)

		res := resourceNsxtPolicyUserManagementRole()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleData())
		d.SetId(roleID)

		err := resourceNsxtPolicyUserManagementRoleDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyUserManagementRole()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleData())

		err := resourceNsxtPolicyUserManagementRoleDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
