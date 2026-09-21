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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	infraMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	lbCookieID          = "lb-cookie-1"
	lbCookieDisplayName = "Test LB Cookie Persistence Profile"
)

func minimalLBCookieData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbCookieDisplayName,
		"nsx_id":       lbCookieID,
	}
}

// lbPersistenceProfileStructValue converts a typed LB persistence profile model into the
// *data.StructValue shape returned by the mocked SDK client's Get/Update calls.
func lbPersistenceProfileStructValue(t *testing.T, obj interface{}, bindingType bindings.BindingType) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(obj, bindingType)
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func setupLBPersistenceProfileMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockLbPersistenceProfilesClient, func()) {
	mockSDK := infraMocks.NewMockLbPersistenceProfilesClient(ctrl)
	mockWrapper := &infraapi.LBPersistenceProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliLbPersistenceProfilesClient
	cliLbPersistenceProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LBPersistenceProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliLbPersistenceProfilesClient = original }
}

func TestMockResourceNsxtPolicyLBCookiePersistenceProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbCookieID).Return(nil, nil)

		res := resourceNsxtPolicyLBCookiePersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBCookieData())

		err := resourceNsxtPolicyLBCookiePersistenceProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create succeeds and re-reads the profile", func(t *testing.T) {
		displayName := lbCookieDisplayName
		sv := lbPersistenceProfileStructValue(t, model.LBCookiePersistenceProfile{
			DisplayName:  &displayName,
			ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE,
		}, model.LBCookiePersistenceProfileBindingType())

		mockSDK.EXPECT().Get(lbCookieID).Return(nil, vapiErrors.NotFound{})
		mockSDK.EXPECT().Patch(lbCookieID, gomock.Any()).Return(nil)
		mockSDK.EXPECT().Get(lbCookieID).Return(sv, nil)

		res := resourceNsxtPolicyLBCookiePersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBCookieData())

		err := resourceNsxtPolicyLBCookiePersistenceProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbCookieID, d.Id())
		assert.Equal(t, lbCookieDisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyLBCookiePersistenceProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbCookieID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBCookiePersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBCookieData())
		d.SetId(lbCookieID)

		err := resourceNsxtPolicyLBCookiePersistenceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbCookieID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLBCookiePersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBCookieData())
		d.SetId(lbCookieID)

		err := resourceNsxtPolicyLBCookiePersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBCookiePersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBCookieData())

		err := resourceNsxtPolicyLBCookiePersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read success populates schema fields", func(t *testing.T) {
		displayName := lbCookieDisplayName
		description := "a cookie profile"
		sv := lbPersistenceProfileStructValue(t, model.LBCookiePersistenceProfile{
			DisplayName:    &displayName,
			Description:    &description,
			ResourceType:   model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE,
			CookieName:     strPtr("NSXLB"),
			CookieMode:     strPtr(model.LBCookiePersistenceProfile_COOKIE_MODE_INSERT),
			CookieFallback: boolPtr(true),
			CookieGarble:   boolPtr(true),
		}, model.LBCookiePersistenceProfileBindingType())
		mockSDK.EXPECT().Get(lbCookieID).Return(sv, nil)

		res := resourceNsxtPolicyLBCookiePersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBCookieData())
		d.SetId(lbCookieID)

		err := resourceNsxtPolicyLBCookiePersistenceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbCookieDisplayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
		assert.Equal(t, "NSXLB", d.Get("cookie_name"))
	})
}

func TestMockResourceNsxtPolicyLBCookiePersistenceProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBCookiePersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBCookieData())

		err := resourceNsxtPolicyLBCookiePersistenceProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update succeeds and re-reads the profile", func(t *testing.T) {
		displayName := lbCookieDisplayName
		sv := lbPersistenceProfileStructValue(t, model.LBCookiePersistenceProfile{
			DisplayName:  &displayName,
			ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE,
		}, model.LBCookiePersistenceProfileBindingType())

		mockSDK.EXPECT().Update(lbCookieID, gomock.Any()).Return(nil, nil)
		mockSDK.EXPECT().Get(lbCookieID).Return(sv, nil)

		res := resourceNsxtPolicyLBCookiePersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBCookieData())
		d.SetId(lbCookieID)

		err := resourceNsxtPolicyLBCookiePersistenceProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyLBCookiePersistenceProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbCookieID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBCookiePersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBCookieData())
		d.SetId(lbCookieID)

		err := resourceNsxtPolicyLBPersistenceProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBCookiePersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBCookieData())

		err := resourceNsxtPolicyLBPersistenceProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
