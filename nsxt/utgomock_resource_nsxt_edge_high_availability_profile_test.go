//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/ClusterProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/ClusterProfilesClient.go ClusterProfilesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"

	nsxmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx"
)

var (
	haProfileID          = "ha-profile-uuid-1"
	haProfileDisplayName = "test-ha-profile"
	haProfileDescription = "HA profile description"
	haProfileRevision    = int64(1)
	haProfileBfdHops     = int64(200)
	haProfileBfdDead     = int64(4)
	haProfileBfdProbe    = int64(600)
	haProfileStandby     = int64(60)
)

// edgeHaProfileStructValue builds the *data.StructValue that the resource's
// Read path expects to receive from the ClusterProfilesClient.
func edgeHaProfileStructValue() *data.StructValue {
	converter := bindings.NewTypeConverter()
	obj := nsxModel.EdgeHighAvailabilityProfile{
		Id:                     &haProfileID,
		DisplayName:            &haProfileDisplayName,
		Description:            &haProfileDescription,
		Revision:               &haProfileRevision,
		BfdAllowedHops:         &haProfileBfdHops,
		BfdDeclareDeadMultiple: &haProfileBfdDead,
		BfdProbeInterval:       &haProfileBfdProbe,
		StandbyRelocationConfig: &nsxModel.StandbyRelocationConfig{
			StandbyRelocationThreshold: &haProfileStandby,
		},
		ResourceType: nsxModel.EdgeHighAvailabilityProfile__TYPE_IDENTIFIER,
	}
	dataValue, errs := converter.ConvertToVapi(obj, nsxModel.EdgeHighAvailabilityProfileBindingType())
	if errs != nil {
		panic(errs[0])
	}
	return dataValue.(*data.StructValue)
}

func setupHAProfileMock(t *testing.T, ctrl *gomock.Controller) (*nsxmocks.MockClusterProfilesClient, func()) {
	mockSDK := nsxmocks.NewMockClusterProfilesClient(ctrl)

	originalCli := cliClusterProfilesClient
	cliClusterProfilesClient = func(_ vapiProtocolClient.Connector) nsx.ClusterProfilesClient {
		return mockSDK
	}
	return mockSDK, func() { cliClusterProfilesClient = originalCli }
}

func haProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":                 haProfileDisplayName,
		"description":                  haProfileDescription,
		"bfd_allowed_hops":             int(haProfileBfdHops),
		"bfd_declare_dead_multiple":    int(haProfileBfdDead),
		"bfd_probe_interval":           int(haProfileBfdProbe),
		"standby_relocation_threshold": int(haProfileStandby),
	}
}

func TestMockResourceNsxtEdgeHighAvailabilityProfileCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHAProfileMock(t, ctrl)
		defer restore()

		sv := edgeHaProfileStructValue()
		mockSDK.EXPECT().Create(gomock.Any()).Return(sv, nil)
		mockSDK.EXPECT().Get(haProfileID).Return(sv, nil)

		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, haProfileData())

		err := resourceNsxtEdgeHighAvailabilityProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, haProfileID, d.Id())
		assert.Equal(t, haProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, int(haProfileBfdHops), d.Get("bfd_allowed_hops"))
	})

	t.Run("Create fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHAProfileMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Create(gomock.Any()).Return(nil, errors.New("create API error"))

		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, haProfileData())

		err := resourceNsxtEdgeHighAvailabilityProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "create API error")
	})
}

func TestMockResourceNsxtEdgeHighAvailabilityProfileRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHAProfileMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(haProfileID).Return(edgeHaProfileStructValue(), nil)

		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(haProfileID)

		err := resourceNsxtEdgeHighAvailabilityProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, haProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, haProfileDescription, d.Get("description"))
		assert.Equal(t, int(haProfileBfdHops), d.Get("bfd_allowed_hops"))
		assert.Equal(t, int(haProfileBfdDead), d.Get("bfd_declare_dead_multiple"))
		assert.Equal(t, int(haProfileBfdProbe), d.Get("bfd_probe_interval"))
		assert.Equal(t, int(haProfileStandby), d.Get("standby_relocation_threshold"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtEdgeHighAvailabilityProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Read fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHAProfileMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(haProfileID).Return(nil, errors.New("read API error"))

		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(haProfileID)

		err := resourceNsxtEdgeHighAvailabilityProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read API error")
	})
}

func TestMockResourceNsxtEdgeHighAvailabilityProfileUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHAProfileMock(t, ctrl)
		defer restore()

		sv := edgeHaProfileStructValue()
		mockSDK.EXPECT().Update(haProfileID, gomock.Any()).Return(sv, nil)
		mockSDK.EXPECT().Get(haProfileID).Return(sv, nil)

		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, haProfileData())
		d.SetId(haProfileID)

		err := resourceNsxtEdgeHighAvailabilityProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, haProfileDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, haProfileData())

		err := resourceNsxtEdgeHighAvailabilityProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Update fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHAProfileMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(haProfileID, gomock.Any()).Return(nil, errors.New("update API error"))

		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, haProfileData())
		d.SetId(haProfileID)

		err := resourceNsxtEdgeHighAvailabilityProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update API error")
	})
}

func TestMockResourceNsxtEdgeHighAvailabilityProfileDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHAProfileMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(haProfileID).Return(nil)

		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(haProfileID)

		err := resourceNsxtEdgeHighAvailabilityProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtEdgeHighAvailabilityProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHAProfileMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(haProfileID).Return(errors.New("delete API error"))

		res := resourceNsxtEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(haProfileID)

		err := resourceNsxtEdgeHighAvailabilityProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete API error")
	})
}
