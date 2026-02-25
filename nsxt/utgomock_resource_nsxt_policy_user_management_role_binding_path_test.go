//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/aaa/RoleBindingsClient.go -package=mocks github.com/vmware/vsphere-automation-sdk-go/services/nsxt/aaa RoleBindingsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/aaa"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	aaamocks "github.com/vmware/terraform-provider-nsxt/mocks/aaa"
)

var (
	rbPathBindingID = "binding-uuid-1"
	rbPathStr       = "/orgs/default"
	rbPathRole      = "auditor"
)

func getNilSlice12() [12]interface{} {
	return [12]interface{}{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}
}

func expectGet(mockSDK *aaamocks.MockRoleBindingsClient, bindingID string, ret nsxModel.RoleBinding, err error) *gomock.Call {
	n := getNilSlice12()
	return mockSDK.EXPECT().Get(bindingID, n[0], n[1], n[2], n[3], n[4], n[5], n[6], n[7], n[8], n[9], n[10], n[11]).Return(ret, err)
}

func minimalRoleBindingPathData() map[string]interface{} {
	return map[string]interface{}{
		"role_binding_id": rbPathBindingID,
		"path":            rbPathStr,
		"roles":           []interface{}{rbPathRole},
	}
}

func roleBindingForPathGetEmpty() nsxModel.RoleBinding {
	return nsxModel.RoleBinding{
		Id:            &rbPathBindingID,
		RolesForPaths: []nsxModel.RolesForPath{},
	}
}

func roleBindingForPathGetWithPath() nsxModel.RoleBinding {
	path := rbPathStr
	role := rbPathRole
	return nsxModel.RoleBinding{
		Id: &rbPathBindingID,
		RolesForPaths: []nsxModel.RolesForPath{
			{
				Path: &path,
				Roles: []nsxModel.Role{
					{Role: &role},
				},
			},
		},
	}
}

func setupRoleBindingsPathMock(t *testing.T, ctrl *gomock.Controller) (*aaamocks.MockRoleBindingsClient, func()) {
	mockSDK := aaamocks.NewMockRoleBindingsClient(ctrl)
	mockWrapper := &aaa.RoleBindingClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliRoleBindingsClient
	cliRoleBindingsClient = func(_ utl.SessionContext, _ client.Connector) *aaa.RoleBindingClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliRoleBindingsClient = original }
}

func TestMockResourceNsxtPolicyUserManagementRoleBindingPathCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRoleBindingsPathMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			expectGet(mockSDK, rbPathBindingID, roleBindingForPathGetEmpty(), nil),
			mockSDK.EXPECT().Update(rbPathBindingID, gomock.Any()).Return(nsxModel.RoleBinding{}, nil),
			expectGet(mockSDK, rbPathBindingID, roleBindingForPathGetWithPath(), nil),
		)

		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleBindingPathData())

		err := resourceNsxtPolicyUserManagementRoleBindingPathCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		wantID := roleBindingPathResourceID(rbPathBindingID, rbPathStr)
		assert.Equal(t, wantID, d.Id())
	})

	t.Run("Create fails when no roles", func(t *testing.T) {
		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		data := minimalRoleBindingPathData()
		data["roles"] = []interface{}{}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyUserManagementRoleBindingPathCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "at least one role")
	})

	t.Run("Create fails when Get returns error", func(t *testing.T) {
		expectGet(mockSDK, rbPathBindingID, nsxModel.RoleBinding{}, errors.New("API error"))

		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleBindingPathData())

		err := resourceNsxtPolicyUserManagementRoleBindingPathCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyUserManagementRoleBindingPathRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRoleBindingsPathMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		expectGet(mockSDK, rbPathBindingID, roleBindingForPathGetWithPath(), nil)

		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(roleBindingPathResourceID(rbPathBindingID, rbPathStr))

		err := resourceNsxtPolicyUserManagementRoleBindingPathRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, rbPathBindingID, d.Get("role_binding_id"))
		assert.Equal(t, rbPathStr, d.Get("path"))
	})

	t.Run("Read fails on invalid resource id", func(t *testing.T) {
		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("not-a-composite-id")

		err := resourceNsxtPolicyUserManagementRoleBindingPathRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role binding path id")
	})

	t.Run("Read clears id when path not on binding", func(t *testing.T) {
		expectGet(mockSDK, rbPathBindingID, roleBindingForPathGetEmpty(), nil)

		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(roleBindingPathResourceID(rbPathBindingID, rbPathStr))

		err := resourceNsxtPolicyUserManagementRoleBindingPathRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Empty(t, d.Id())
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Read clears id when binding not found", func(t *testing.T) {
		expectGet(mockSDK, rbPathBindingID, nsxModel.RoleBinding{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(roleBindingPathResourceID(rbPathBindingID, rbPathStr))

		err := resourceNsxtPolicyUserManagementRoleBindingPathRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyUserManagementRoleBindingPathUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRoleBindingsPathMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			expectGet(mockSDK, rbPathBindingID, roleBindingForPathGetWithPath(), nil),
			mockSDK.EXPECT().Update(rbPathBindingID, gomock.Any()).Return(nsxModel.RoleBinding{}, nil),
			expectGet(mockSDK, rbPathBindingID, roleBindingForPathGetWithPath(), nil),
		)

		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleBindingPathData())
		d.SetId(roleBindingPathResourceID(rbPathBindingID, rbPathStr))

		err := resourceNsxtPolicyUserManagementRoleBindingPathUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when no roles", func(t *testing.T) {
		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		data := minimalRoleBindingPathData()
		data["roles"] = []interface{}{}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(roleBindingPathResourceID(rbPathBindingID, rbPathStr))

		err := resourceNsxtPolicyUserManagementRoleBindingPathUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "at least one role")
	})
}

func TestMockResourceNsxtPolicyUserManagementRoleBindingPathDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRoleBindingsPathMock(t, ctrl)
	defer restore()

	t.Run("Delete success removes path", func(t *testing.T) {
		gomock.InOrder(
			expectGet(mockSDK, rbPathBindingID, roleBindingForPathGetWithPath(), nil),
			mockSDK.EXPECT().Update(rbPathBindingID, gomock.Any()).Return(nsxModel.RoleBinding{}, nil),
		)

		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleBindingPathData())
		d.SetId(roleBindingPathResourceID(rbPathBindingID, rbPathStr))

		err := resourceNsxtPolicyUserManagementRoleBindingPathDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete no-op when path already absent", func(t *testing.T) {
		expectGet(mockSDK, rbPathBindingID, roleBindingForPathGetEmpty(), nil)

		res := resourceNsxtPolicyUserManagementRoleBindingPath()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRoleBindingPathData())
		d.SetId(roleBindingPathResourceID(rbPathBindingID, rbPathStr))

		err := resourceNsxtPolicyUserManagementRoleBindingPathDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
