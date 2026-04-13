//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/aaa/RoleBindingsClient.go -package=mocks github.com/vmware/vsphere-automation-sdk-go/services/nsxt/aaa RoleBindingsClient

package nsxt

import (
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

func setupRoleBindingsDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*aaamocks.MockRoleBindingsClient, func()) {
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

func expectRoleBindingGet(mockSDK *aaamocks.MockRoleBindingsClient, id string, ret nsxModel.RoleBinding, err error) *gomock.Call {
	n := getNilSlice12()
	return mockSDK.EXPECT().Get(id, n[0], n[1], n[2], n[3], n[4], n[5], n[6], n[7], n[8], n[9], n[10], n[11]).Return(ret, err)
}

func dsRoleBindingAPI(id, name, userType string) nsxModel.RoleBinding {
	path := "/"
	role := "auditor"
	return nsxModel.RoleBinding{
		Id:    &id,
		Name:  &name,
		Type_: &userType,
		RolesForPaths: []nsxModel.RolesForPath{
			{Path: &path, Roles: []nsxModel.Role{{Role: &role}}},
		},
	}
}

func TestMockDataSourceNsxtPolicyUserManagementRoleBindingReadByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRoleBindingsDataSourceMock(t, ctrl)
	defer restore()

	t.Run("Read by id success", func(t *testing.T) {
		id := "rb-ds-1"
		name := "user1"
		rbType := nsxModel.RoleBinding_TYPE_REMOTE_USER
		expectRoleBindingGet(mockSDK, id, dsRoleBindingAPI(id, name, rbType), nil)

		ds := dataSourceNsxtPolicyRoleBinding()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": id,
		})

		err := dataSourceNsxtPolicyRoleBindingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, id, d.Id())
		assert.Equal(t, name, d.Get("name"))
		assert.Equal(t, rbType, d.Get("type"))
	})

	t.Run("Read by id API error", func(t *testing.T) {
		id := "rb-missing"
		expectRoleBindingGet(mockSDK, id, nsxModel.RoleBinding{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyRoleBinding()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": id,
		})

		err := dataSourceNsxtPolicyRoleBindingRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockDataSourceNsxtPolicyUserManagementRoleBindingReadByNameType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRoleBindingsDataSourceMock(t, ctrl)
	defer restore()

	t.Run("Read by name and type success", func(t *testing.T) {
		id := "rb-ds-2"
		name := "remoteuser"
		rbType := nsxModel.RoleBinding_TYPE_REMOTE_USER
		userTypeParam := rbType
		mockSDK.EXPECT().List(nil, nil, nil, nil, &name, nil, nil, nil, nil, nil, nil, &userTypeParam).Return(nsxModel.RoleBindingListResult{
			Results: []nsxModel.RoleBinding{dsRoleBindingAPI(id, name, rbType)},
		}, nil)

		ds := dataSourceNsxtPolicyRoleBinding()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"name": name,
			"type": rbType,
		})

		err := dataSourceNsxtPolicyRoleBindingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, id, d.Id())
	})

	t.Run("Read by name and type not found", func(t *testing.T) {
		name := "nouser"
		rbType := nsxModel.RoleBinding_TYPE_REMOTE_USER
		userTypeParam := rbType
		mockSDK.EXPECT().List(nil, nil, nil, nil, &name, nil, nil, nil, nil, nil, nil, &userTypeParam).Return(nsxModel.RoleBindingListResult{
			Results: []nsxModel.RoleBinding{},
		}, nil)

		ds := dataSourceNsxtPolicyRoleBinding()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"name": name,
			"type": rbType,
		})

		err := dataSourceNsxtPolicyRoleBindingRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Read fails when neither id nor name+type", func(t *testing.T) {
		ds := dataSourceNsxtPolicyRoleBinding()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyRoleBindingRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "either id or both name and type")
	})
}
