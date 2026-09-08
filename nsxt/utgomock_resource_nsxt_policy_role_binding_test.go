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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

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

func TestUnitNsxt_rolesPerPath_getAnyRole(t *testing.T) {
	t.Run("nil map returns nil", func(t *testing.T) {
		var r rolesPerPath
		assert.Nil(t, r.getAnyRole())
	})

	t.Run("no true values returns nil", func(t *testing.T) {
		r := rolesPerPath{"auditor": false}
		assert.Nil(t, r.getAnyRole())
	})

	t.Run("returns the true role", func(t *testing.T) {
		r := rolesPerPath{"auditor": true}
		role := r.getAnyRole()
		require.NotNil(t, role)
		assert.Equal(t, "auditor", *role)
	})
}

func TestUnitNsxt_getRolesForPathFromSchema(t *testing.T) {
	res := resourceNsxtPolicyUserManagementRoleBinding()

	t.Run("empty when roles_for_path is unset", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		got := getRolesForPathFromSchema(d)
		assert.Empty(t, got)
	})

	t.Run("parses a roles_for_path set", func(t *testing.T) {
		data := minimalRbData()
		data["roles_for_path"] = []interface{}{
			map[string]interface{}{
				"path":  "/",
				"roles": []interface{}{"auditor"},
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		got := getRolesForPathFromSchema(d)
		require.Contains(t, got, "/")
		assert.True(t, got["/"]["auditor"])
	})
}

func TestUnitNsxt_getRolesForPathList(t *testing.T) {
	res := resourceNsxtPolicyUserManagementRoleBinding()

	t.Run("converts current roles_for_path with no removals", func(t *testing.T) {
		data := minimalRbData()
		data["roles_for_path"] = []interface{}{
			map[string]interface{}{
				"path":  "/",
				"roles": []interface{}{"auditor"},
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		got := getRolesForPathList(d, rolesForPath{})
		require.Len(t, got, 1)
		assert.Equal(t, "/", *got[0].Path)
	})

	t.Run("appends a DeletePath entry for a path being removed", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		toRemove := rolesForPath{"/other-path": rolesPerPath{"auditor": true}}
		got := getRolesForPathList(d, toRemove)
		require.Len(t, got, 1)
		assert.Equal(t, "/other-path", *got[0].Path)
		require.NotNil(t, got[0].DeletePath)
		assert.True(t, *got[0].DeletePath)
	})

	t.Run("skips removal for a path also present in the current definition", func(t *testing.T) {
		data := minimalRbData()
		data["roles_for_path"] = []interface{}{
			map[string]interface{}{
				"path":  "/",
				"roles": []interface{}{"auditor"},
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		toRemove := rolesForPath{"/": rolesPerPath{"auditor": true}}
		got := getRolesForPathList(d, toRemove)
		require.Len(t, got, 1)
		assert.Equal(t, "/", *got[0].Path)
		assert.Nil(t, got[0].DeletePath)
	})
}

func TestUnitNsxt_setRolesForPathInSchema(t *testing.T) {
	res := resourceNsxtPolicyUserManagementRoleBinding()

	t.Run("sets roles_for_path from the API response", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		path := "/"
		role := "auditor"
		setRolesForPathInSchema(d, []nsxModel.RolesForPath{
			{Path: &path, Roles: []nsxModel.Role{{Role: &role}}},
		})
		got := d.Get("roles_for_path").(*schema.Set).List()
		require.Len(t, got, 1)
		elem := got[0].(map[string]interface{})
		assert.Equal(t, "/", elem["path"])
	})

	t.Run("skips entries with a nil Path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		setRolesForPathInSchema(d, []nsxModel.RolesForPath{{Path: nil}})
		got := d.Get("roles_for_path").(*schema.Set).List()
		assert.Empty(t, got)
	})
}

func TestMockNsxt_getExistingRoleBinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRbMock(t, ctrl)
	defer restore()
	rbClient := &aaaapi.RoleBindingClientContext{Client: mockSDK, ClientType: utl.Local}

	t.Run("finds a matching binding", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, &rbName, nil, nil, nil, nil, nil, nil, nil).Return(
			nsxModel.RoleBindingListResult{Results: []nsxModel.RoleBinding{rbAPIResponse()}}, nil,
		)
		obj, err := getExistingRoleBinding(rbClient, rbName, rbType)
		require.NoError(t, err)
		assert.Equal(t, rbID, *obj.Id)
	})

	t.Run("List error propagates", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, &rbName, nil, nil, nil, nil, nil, nil, nil).Return(
			nsxModel.RoleBindingListResult{}, vapiErrors.InternalServerError{},
		)
		_, err := getExistingRoleBinding(rbClient, rbName, rbType)
		require.Error(t, err)
	})

	t.Run("no matching name/type returns an error", func(t *testing.T) {
		other := "someone-else"
		mockSDK.EXPECT().List(nil, nil, nil, nil, &rbName, nil, nil, nil, nil, nil, nil, nil).Return(
			nsxModel.RoleBindingListResult{Results: []nsxModel.RoleBinding{{Name: &other, Type_: &rbType}}}, nil,
		)
		_, err := getExistingRoleBinding(rbClient, rbName, rbType)
		require.Error(t, err)
	})
}

func TestMockNsxt_overwriteRoleBinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRbMock(t, ctrl)
	defer restore()

	t.Run("success updates then reads", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(rbID, gomock.Any()).Return(rbAPIResponse(), nil),
			mockSDK.EXPECT().Get(rbID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(rbAPIResponse(), nil),
		)
		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		d.SetId(rbID)

		localType := "local_user"
		existing := rbAPIResponse()
		path := "/"
		existing.RolesForPaths = []nsxModel.RolesForPath{{Path: &path, Roles: []nsxModel.Role{{Role: &localType}}}}

		err := overwriteRoleBinding(d, newGoMockProviderClient(), &existing)
		require.NoError(t, err)
	})

	t.Run("Update error propagates", func(t *testing.T) {
		mockSDK.EXPECT().Update(rbID, gomock.Any()).Return(nsxModel.RoleBinding{}, vapiErrors.InternalServerError{})
		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		d.SetId(rbID)

		err := overwriteRoleBinding(d, newGoMockProviderClient(), &nsxModel.RoleBinding{})
		require.Error(t, err)
	})
}

func TestMockNsxt_revertRoleBinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRbMock(t, ctrl)
	defer restore()

	t.Run("success reverts to auditor role", func(t *testing.T) {
		mockSDK.EXPECT().Update(rbID, gomock.Any()).Return(rbAPIResponse(), nil)
		res := resourceNsxtPolicyUserManagementRoleBinding()
		data := minimalRbData()
		data["roles_for_path"] = []interface{}{
			map[string]interface{}{
				"path":  "/some/other/path",
				"roles": []interface{}{"auditor"},
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(rbID)

		err := revertRoleBinding(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update error propagates", func(t *testing.T) {
		mockSDK.EXPECT().Update(rbID, gomock.Any()).Return(nsxModel.RoleBinding{}, vapiErrors.InternalServerError{})
		res := resourceNsxtPolicyUserManagementRoleBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRbData())
		d.SetId(rbID)

		err := revertRoleBinding(d, newGoMockProviderClient())
		require.Error(t, err)
	})
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
