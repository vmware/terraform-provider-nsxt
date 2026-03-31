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

	t0fpapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0fpmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
)

var (
	fppBindingID          = "default"
	fppBindingDisplayName = "Test FPP Binding"
	fppBindingDescription = "Test flood protection profile binding"
	fppBindingRevision    = int64(1)
	fppBindingT0ID        = "t0-fp-gw-1"
	fppBindingParentPath  = "/infra/tier-0s/t0-fp-gw-1"
	fppBindingPath        = "/infra/tier-0s/t0-fp-gw-1/flood-protection-profile-bindings/default"
	fppBindingProfilePath = "/infra/flood-protection-profiles/gw-fpp-1"
)

func fppBindingAPIResponse() nsxModel.FloodProtectionProfileBindingMap {
	return nsxModel.FloodProtectionProfileBindingMap{
		Id:          &fppBindingID,
		DisplayName: &fppBindingDisplayName,
		Description: &fppBindingDescription,
		Revision:    &fppBindingRevision,
		Path:        &fppBindingPath,
		ProfilePath: &fppBindingProfilePath,
	}
}

func minimalFppBindingData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": fppBindingDisplayName,
		"description":  fppBindingDescription,
		"parent_path":  fppBindingParentPath,
		"profile_path": fppBindingProfilePath,
	}
}

func setupFppBindingMock(t *testing.T, ctrl *gomock.Controller) (*t0fpmocks.MockFloodProtectionProfileBindingsClient, func()) {
	mockSDK := t0fpmocks.NewMockFloodProtectionProfileBindingsClient(ctrl)
	mockWrapper := &t0fpapi.FloodProtectionProfileBindingMapClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliTier0FloodProtectionProfileBindingsClient
	cliTier0FloodProtectionProfileBindingsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t0fpapi.FloodProtectionProfileBindingMapClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier0FloodProtectionProfileBindingsClient = original }
}

func TestMockResourceNsxtPolicyGatewayFloodProtectionProfileBindingCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupFppBindingMock(t, ctrl)
	defer restore()

	t.Run("Create success for Tier0 parent path", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(fppBindingT0ID, fppBindingID).Return(nsxModel.FloodProtectionProfileBindingMap{}, notFoundErr),
			mockSDK.EXPECT().Patch(fppBindingT0ID, fppBindingID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(fppBindingT0ID, fppBindingID).Return(fppBindingAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFppBindingData())

		err := resourceNsxtPolicyGatewayFloodProtectionProfileBindingCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, fppBindingID, d.Id())
	})

	t.Run("Create fails when binding already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(fppBindingT0ID, fppBindingID).Return(fppBindingAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFppBindingData())

		err := resourceNsxtPolicyGatewayFloodProtectionProfileBindingCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayFloodProtectionProfileBindingRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupFppBindingMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(fppBindingT0ID, fppBindingID).Return(fppBindingAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFppBindingData())
		d.SetId(fppBindingID)

		err := resourceNsxtPolicyGatewayFloodProtectionProfileBindingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, fppBindingDisplayName, d.Get("display_name"))
		assert.Equal(t, fppBindingProfilePath, d.Get("profile_path"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(fppBindingT0ID, fppBindingID).Return(nsxModel.FloodProtectionProfileBindingMap{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFppBindingData())
		d.SetId(fppBindingID)

		err := resourceNsxtPolicyGatewayFloodProtectionProfileBindingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFppBindingData())

		err := resourceNsxtPolicyGatewayFloodProtectionProfileBindingRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayFloodProtectionProfileBindingUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupFppBindingMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(fppBindingT0ID, fppBindingID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(fppBindingT0ID, fppBindingID).Return(fppBindingAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFppBindingData())
		d.SetId(fppBindingID)

		err := resourceNsxtPolicyGatewayFloodProtectionProfileBindingUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFppBindingData())

		err := resourceNsxtPolicyGatewayFloodProtectionProfileBindingUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayFloodProtectionProfileBindingDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupFppBindingMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(fppBindingT0ID, fppBindingID).Return(nil)

		res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFppBindingData())
		d.SetId(fppBindingID)

		err := resourceNsxtPolicyGatewayFloodProtectionProfileBindingDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFppBindingData())

		err := resourceNsxtPolicyGatewayFloodProtectionProfileBindingDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
