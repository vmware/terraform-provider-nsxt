// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/HostSwitchProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/HostSwitchProfilesClient.go HostSwitchProfilesClient

package nsxt

import (
	"errors"
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
	vtepHAProfileID          = "vtep-ha-profile-1"
	vtepHAProfileDisplayName = "test-vtep-ha-profile"
	vtepHAProfileDescription = "VtepHA host switch profile description"
	vtepHAProfilePath        = "/infra/host-switch-profiles/vtep-ha-profile-1"
	vtepHAProfileRevision    = int64(1)
)

func vtepHAProfileStructValue() *data.StructValue {
	autoRecovery := true
	autoRecoveryInitialWait := int64(300)
	autoRecoveryMaxBackoff := int64(86400)
	enabled := false
	failoverTimeout := int64(5)
	profile := model.PolicyVtepHAHostSwitchProfile{
		Id:                      &vtepHAProfileID,
		DisplayName:             &vtepHAProfileDisplayName,
		Description:             &vtepHAProfileDescription,
		Path:                    &vtepHAProfilePath,
		Revision:                &vtepHAProfileRevision,
		ResourceType:            model.PolicyBaseHostSwitchProfile_RESOURCE_TYPE_POLICYVTEPHAHOSTSWITCHPROFILE,
		AutoRecovery:            &autoRecovery,
		AutoRecoveryInitialWait: &autoRecoveryInitialWait,
		AutoRecoveryMaxBackoff:  &autoRecoveryMaxBackoff,
		Enabled:                 &enabled,
		FailoverTimeout:         &failoverTimeout,
	}
	converter := bindings.NewTypeConverter()
	dataValue, errs := converter.ConvertToVapi(profile, model.PolicyVtepHAHostSwitchProfileBindingType())
	if errs != nil {
		panic(errs[0])
	}
	return dataValue.(*data.StructValue)
}

func minimalVtepHAProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vtepHAProfileDisplayName,
		"description":  vtepHAProfileDescription,
	}
}

func TestMockResourceNsxtPolicyVtepHAHostSwitchProfileCreate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK := inframocks.NewMockHostSwitchProfilesClient(ctrl)
	mockWrapper := &cliinfra.HostSwitchProfilesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliHostSwitchProfilesClient
	defer func() { cliHostSwitchProfilesClient = originalCli }()
	cliHostSwitchProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.HostSwitchProfilesClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockSDK.EXPECT().Get(gomock.Any()).Return(vtepHAProfileStructValue(), nil)

		res := resourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVtepHAProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, vtepHAProfileDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(vtepHAProfileID).Return(vtepHAProfileStructValue(), nil)

		res := resourceNsxtVtepHAHostSwitchProfile()
		data := minimalVtepHAProfileData()
		data["nsx_id"] = vtepHAProfileID
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when Patch returns error", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil, errors.New("API error"))

		res := resourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVtepHAProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtPolicyVtepHAHostSwitchProfileRead(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK := inframocks.NewMockHostSwitchProfilesClient(ctrl)
	mockWrapper := &cliinfra.HostSwitchProfilesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliHostSwitchProfilesClient
	defer func() { cliHostSwitchProfilesClient = originalCli }()
	cliHostSwitchProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.HostSwitchProfilesClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(vtepHAProfileID).Return(vtepHAProfileStructValue(), nil)

		res := resourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vtepHAProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, vtepHAProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, vtepHAProfileDescription, d.Get("description"))
		assert.Equal(t, vtepHAProfilePath, d.Get("path"))
		assert.Equal(t, int(vtepHAProfileRevision), d.Get("revision"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining VtepHAHostSwitchProfile ID")
	})

	t.Run("Read clears ID when not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(vtepHAProfileID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vtepHAProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyVtepHAHostSwitchProfileUpdate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK := inframocks.NewMockHostSwitchProfilesClient(ctrl)
	mockWrapper := &cliinfra.HostSwitchProfilesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliHostSwitchProfilesClient
	defer func() { cliHostSwitchProfilesClient = originalCli }()
	cliHostSwitchProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.HostSwitchProfilesClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Update(vtepHAProfileID, gomock.Any()).Return(nil, nil)
		mockSDK.EXPECT().Get(vtepHAProfileID).Return(vtepHAProfileStructValue(), nil)

		res := resourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVtepHAProfileData())
		d.SetId(vtepHAProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, vtepHAProfileDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVtepHAProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining VtepHAHostSwitchProfile ID")
	})
}

func TestMockResourceNsxtPolicyVtepHAHostSwitchProfileDelete(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK := inframocks.NewMockHostSwitchProfilesClient(ctrl)
	mockWrapper := &cliinfra.HostSwitchProfilesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliHostSwitchProfilesClient
	defer func() { cliHostSwitchProfilesClient = originalCli }()
	cliHostSwitchProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.HostSwitchProfilesClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(vtepHAProfileID).Return(nil)

		res := resourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vtepHAProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining VtepHAHostSwitchProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockSDK.EXPECT().Delete(vtepHAProfileID).Return(errors.New("API error"))

		res := resourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vtepHAProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtVtepHAHostSwitchProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}
