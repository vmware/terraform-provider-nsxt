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

	t0fpapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	t0lsfpapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services"
	t1fpapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s"
	t1lsfpapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/locale_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0fpmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	t0lsfpmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services"
	t1fpmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
	t1lsfpmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/locale_services"
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

func TestMockResourceNsxtPolicyGatewayFloodProtectionProfileBindingOtherParents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockT0LS := t0lsfpmocks.NewMockFloodProtectionProfileBindingsClient(ctrl)
	mockT1 := t1fpmocks.NewMockFloodProtectionProfileBindingsClient(ctrl)
	mockT1LS := t1lsfpmocks.NewMockFloodProtectionProfileBindingsClient(ctrl)

	origT0LS := cliT0LocaleServicesFloodProtectionProfileBindingsClient
	origT1 := cliTier1FloodProtectionProfileBindingsClient
	origT1LS := cliT1LocaleServicesFloodProtectionProfileBindingsClient
	defer func() {
		cliT0LocaleServicesFloodProtectionProfileBindingsClient = origT0LS
		cliTier1FloodProtectionProfileBindingsClient = origT1
		cliT1LocaleServicesFloodProtectionProfileBindingsClient = origT1LS
	}()
	cliT0LocaleServicesFloodProtectionProfileBindingsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t0lsfpapi.FloodProtectionProfileBindingMapClientContext {
		return &t0lsfpapi.FloodProtectionProfileBindingMapClientContext{Client: mockT0LS, ClientType: utl.Local}
	}
	cliTier1FloodProtectionProfileBindingsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t1fpapi.FloodProtectionProfileBindingMapClientContext {
		return &t1fpapi.FloodProtectionProfileBindingMapClientContext{Client: mockT1, ClientType: utl.Local}
	}
	cliT1LocaleServicesFloodProtectionProfileBindingsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t1lsfpapi.FloodProtectionProfileBindingMapClientContext {
		return &t1lsfpapi.FloodProtectionProfileBindingMapClientContext{Client: mockT1LS, ClientType: utl.Local}
	}

	res := resourceNsxtPolicyGatewayFloodProtectionProfileBinding()
	notFound := vapiErrors.NotFound{}

	run := func(t *testing.T, parentPath string) {
		data := minimalFppBindingData()
		data["parent_path"] = parentPath
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		require.NoError(t, resourceNsxtPolicyGatewayFloodProtectionProfileBindingCreate(d, newGoMockProviderClient()))
		assert.Equal(t, fppBindingID, d.Id())
		require.NoError(t, resourceNsxtPolicyGatewayFloodProtectionProfileBindingDelete(d, newGoMockProviderClient()))
	}

	t.Run("Tier0 locale service parent", func(t *testing.T) {
		gomock.InOrder(
			mockT0LS.EXPECT().Get("t0", "ls", fppBindingID).Return(nsxModel.FloodProtectionProfileBindingMap{}, notFound),
			mockT0LS.EXPECT().Patch("t0", "ls", fppBindingID, gomock.Any()).Return(nil),
			mockT0LS.EXPECT().Get("t0", "ls", fppBindingID).Return(fppBindingAPIResponse(), nil),
			mockT0LS.EXPECT().Delete("t0", "ls", fppBindingID).Return(nil),
		)
		run(t, "/infra/tier-0s/t0/locale-services/ls")
	})

	t.Run("Tier1 parent", func(t *testing.T) {
		gomock.InOrder(
			mockT1.EXPECT().Get("t1", fppBindingID).Return(nsxModel.FloodProtectionProfileBindingMap{}, notFound),
			mockT1.EXPECT().Patch("t1", fppBindingID, gomock.Any()).Return(nil),
			mockT1.EXPECT().Get("t1", fppBindingID).Return(fppBindingAPIResponse(), nil),
			mockT1.EXPECT().Delete("t1", fppBindingID).Return(nil),
		)
		run(t, "/infra/tier-1s/t1")
	})

	t.Run("Tier1 locale service parent", func(t *testing.T) {
		gomock.InOrder(
			mockT1LS.EXPECT().Get("t1", "ls", fppBindingID).Return(nsxModel.FloodProtectionProfileBindingMap{}, notFound),
			mockT1LS.EXPECT().Patch("t1", "ls", fppBindingID, gomock.Any()).Return(nil),
			mockT1LS.EXPECT().Get("t1", "ls", fppBindingID).Return(fppBindingAPIResponse(), nil),
			mockT1LS.EXPECT().Delete("t1", "ls", fppBindingID).Return(nil),
		)
		run(t, "/infra/tier-1s/t1/locale-services/ls")
	})

	t.Run("invalid parent path fails", func(t *testing.T) {
		_, err := resourceNsxtPolicyGatewayFloodProtectionProfileBindingGet(utl.SessionContext{ClientType: utl.Local}, nil, "/infra/segments/s1", fppBindingID)
		require.Error(t, err)
	})
}
