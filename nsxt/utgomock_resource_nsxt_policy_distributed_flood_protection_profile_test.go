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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dfppID       = "dfpp-1"
	dfppName     = "dfpp-fooname"
	dfppRevision = int64(1)
	dfppPath     = "/infra/flood-protection-profiles/dfpp-1"
)

func dfppStructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	obj := model.DistributedFloodProtectionProfile{
		DisplayName:  &dfppName,
		Path:         &dfppPath,
		Revision:     &dfppRevision,
		ResourceType: model.FloodProtectionProfile_RESOURCE_TYPE_DISTRIBUTEDFLOODPROTECTIONPROFILE,
	}
	val, errs := converter.ConvertToVapi(obj, model.DistributedFloodProtectionProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestMockResourceNsxtPolicyDistributedFloodProtectionProfileCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDfppSDK := inframocks.NewMockFloodProtectionProfilesClient(ctrl)
	dfppWrapper := &cliinfra.StructValueClientContext{
		Client:     mockDfppSDK,
		ClientType: utl.Local,
	}

	originalCli := cliFloodProtectionProfilesClient
	defer func() { cliFloodProtectionProfilesClient = originalCli }()
	cliFloodProtectionProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.StructValueClientContext {
		return dfppWrapper
	}

	res := resourceNsxtPolicyDistributedFloodProtectionProfile()

	t.Run("Create_success", func(t *testing.T) {
		mockDfppSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		sv := dfppStructValue(t)
		mockDfppSDK.EXPECT().Get(gomock.Any()).Return(sv, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": dfppName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, dfppName, d.Get("display_name"))
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockDfppSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": dfppName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyDistributedFloodProtectionProfileRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDfppSDK := inframocks.NewMockFloodProtectionProfilesClient(ctrl)
	dfppWrapper := &cliinfra.StructValueClientContext{
		Client:     mockDfppSDK,
		ClientType: utl.Local,
	}

	originalCli := cliFloodProtectionProfilesClient
	defer func() { cliFloodProtectionProfilesClient = originalCli }()
	cliFloodProtectionProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.StructValueClientContext {
		return dfppWrapper
	}

	res := resourceNsxtPolicyDistributedFloodProtectionProfile()

	t.Run("Read_success", func(t *testing.T) {
		sv := dfppStructValue(t)
		mockDfppSDK.EXPECT().Get(dfppID).Return(sv, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dfppID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, dfppName, d.Get("display_name"))
	})

	t.Run("Read_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FloodProtectionProfile ID")
	})

	t.Run("Read_fails_when_Get_returns_not_found", func(t *testing.T) {
		mockDfppSDK.EXPECT().Get(dfppID).Return(nil, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dfppID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyDistributedFloodProtectionProfileUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDfppSDK := inframocks.NewMockFloodProtectionProfilesClient(ctrl)
	dfppWrapper := &cliinfra.StructValueClientContext{
		Client:     mockDfppSDK,
		ClientType: utl.Local,
	}

	originalCli := cliFloodProtectionProfilesClient
	defer func() { cliFloodProtectionProfilesClient = originalCli }()
	cliFloodProtectionProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.StructValueClientContext {
		return dfppWrapper
	}

	res := resourceNsxtPolicyDistributedFloodProtectionProfile()

	t.Run("Update_success", func(t *testing.T) {
		mockDfppSDK.EXPECT().Patch(dfppID, gomock.Any(), gomock.Any()).Return(nil)
		sv := dfppStructValue(t)
		mockDfppSDK.EXPECT().Get(dfppID).Return(sv, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": dfppName,
		})
		d.SetId(dfppID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyDistributedFloodProtectionProfileDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDfppSDK := inframocks.NewMockFloodProtectionProfilesClient(ctrl)
	dfppWrapper := &cliinfra.StructValueClientContext{
		Client:     mockDfppSDK,
		ClientType: utl.Local,
	}

	originalCli := cliFloodProtectionProfilesClient
	defer func() { cliFloodProtectionProfilesClient = originalCli }()
	cliFloodProtectionProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.StructValueClientContext {
		return dfppWrapper
	}

	res := resourceNsxtPolicyDistributedFloodProtectionProfile()

	t.Run("Delete_success", func(t *testing.T) {
		mockDfppSDK.EXPECT().Delete(dfppID, gomock.Any()).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dfppID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyFloodProtectionProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyFloodProtectionProfileDelete(d, m)
		require.Error(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockDfppSDK.EXPECT().Delete(dfppID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dfppID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyFloodProtectionProfileDelete(d, m)
		require.Error(t, err)
	})
}
