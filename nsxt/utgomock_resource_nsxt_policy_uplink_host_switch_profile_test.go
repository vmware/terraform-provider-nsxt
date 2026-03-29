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
	uplinkProfileID          = "uplink-profile-1"
	uplinkProfileDisplayName = "test-uplink-profile"
	uplinkProfileDescription = "Uplink host switch profile description"
	uplinkProfilePath        = "/infra/host-switch-profiles/uplink-profile-1"
	uplinkProfileRevision    = int64(1)
)

func uplinkProfileStructValue() *data.StructValue {
	uplinkName := "uplink1"
	uplinkType := model.Uplink_UPLINK_TYPE_PNIC
	policy := model.TeamingPolicy_POLICY_FAILOVER_ORDER
	geneve := model.PolicyUplinkHostSwitchProfile_OVERLAY_ENCAP_GENEVE
	transportVlan := int64(0)
	profile := model.PolicyUplinkHostSwitchProfile{
		Id:           &uplinkProfileID,
		DisplayName:  &uplinkProfileDisplayName,
		Description:  &uplinkProfileDescription,
		Path:         &uplinkProfilePath,
		Revision:     &uplinkProfileRevision,
		ResourceType: model.PolicyBaseHostSwitchProfile_RESOURCE_TYPE_POLICYUPLINKHOSTSWITCHPROFILE,
		Teaming: &model.TeamingPolicy{
			Policy:     &policy,
			ActiveList: []model.Uplink{{UplinkName: &uplinkName, UplinkType: &uplinkType}},
		},
		OverlayEncap:  &geneve,
		TransportVlan: &transportVlan,
	}
	converter := bindings.NewTypeConverter()
	dataValue, errs := converter.ConvertToVapi(profile, model.PolicyUplinkHostSwitchProfileBindingType())
	if errs != nil {
		panic(errs[0])
	}
	return dataValue.(*data.StructValue)
}

func minimalUplinkProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": uplinkProfileDisplayName,
		"description":  uplinkProfileDescription,
		"teaming": []interface{}{
			map[string]interface{}{
				"active": []interface{}{
					map[string]interface{}{
						"uplink_name": "uplink1",
						"uplink_type": model.Uplink_UPLINK_TYPE_PNIC,
					},
				},
				"policy":  model.TeamingPolicy_POLICY_FAILOVER_ORDER,
				"standby": []interface{}{},
			},
		},
	}
}

func TestMockResourceNsxtPolicyUplinkHostSwitchProfileCreate(t *testing.T) {
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
		mockSDK.EXPECT().Get(gomock.Any()).Return(uplinkProfileStructValue(), nil)

		res := resourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUplinkProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, uplinkProfileDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(uplinkProfileID).Return(uplinkProfileStructValue(), nil)

		res := resourceNsxtUplinkHostSwitchProfile()
		data := minimalUplinkProfileData()
		data["nsx_id"] = uplinkProfileID
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when Patch returns error", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil, errors.New("API error"))

		res := resourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUplinkProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtPolicyUplinkHostSwitchProfileRead(t *testing.T) {
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
		mockSDK.EXPECT().Get(uplinkProfileID).Return(uplinkProfileStructValue(), nil)

		res := resourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(uplinkProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, uplinkProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, uplinkProfileDescription, d.Get("description"))
		assert.Equal(t, uplinkProfilePath, d.Get("path"))
		assert.Equal(t, int(uplinkProfileRevision), d.Get("revision"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining UplinkHostSwitchProfile ID")
	})

	t.Run("Read clears ID when not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(uplinkProfileID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(uplinkProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyUplinkHostSwitchProfileUpdate(t *testing.T) {
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
		mockSDK.EXPECT().Update(uplinkProfileID, gomock.Any()).Return(nil, nil)
		mockSDK.EXPECT().Get(uplinkProfileID).Return(uplinkProfileStructValue(), nil)

		res := resourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUplinkProfileData())
		d.SetId(uplinkProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, uplinkProfileDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUplinkProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining UplinkHostSwitchProfile ID")
	})
}

func TestMockResourceNsxtPolicyUplinkHostSwitchProfileDelete(t *testing.T) {
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
		mockSDK.EXPECT().Delete(uplinkProfileID).Return(nil)

		res := resourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(uplinkProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining UplinkHostSwitchProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockSDK.EXPECT().Delete(uplinkProfileID).Return(errors.New("API error"))

		res := resourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(uplinkProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtUplinkHostSwitchProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}
