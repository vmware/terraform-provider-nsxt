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

	tier0sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	intervrfmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
)

var (
	interVrfID          = "inter-vrf-1"
	interVrfDisplayName = "Test InterVRF Route"
	interVrfDescription = "Test inter VRF routing"
	interVrfRevision    = int64(1)
	interVrfGwPath      = "/infra/tier-0s/t0-gw-1"
	interVrfGwID        = "t0-gw-1"
	interVrfTargetPath  = "/infra/tier-0s/t0-gw-2"
)

func interVrfAPIResponse() nsxModel.PolicyInterVrfRoutingConfig {
	return nsxModel.PolicyInterVrfRoutingConfig{
		Id:          &interVrfID,
		DisplayName: &interVrfDisplayName,
		Description: &interVrfDescription,
		Revision:    &interVrfRevision,
		TargetPath:  &interVrfTargetPath,
	}
}

func minimalInterVrfData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": interVrfDisplayName,
		"description":  interVrfDescription,
		"nsx_id":       interVrfID,
		"gateway_path": interVrfGwPath,
		"target_path":  interVrfTargetPath,
	}
}

func setupInterVrfMock(t *testing.T, ctrl *gomock.Controller) (*intervrfmocks.MockInterVrfRoutingClient, func()) {
	mockSDK := intervrfmocks.NewMockInterVrfRoutingClient(ctrl)
	mockWrapper := &tier0sapi.PolicyInterVrfRoutingConfigClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliInterVrfRoutingClient
	cliInterVrfRoutingClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0sapi.PolicyInterVrfRoutingConfigClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliInterVrfRoutingClient = original }
}

func TestMockResourceNsxtPolicyTier0InterVRFRoutingCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupInterVrfMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(nsxModel.PolicyInterVrfRoutingConfig{}, notFoundErr),
			mockSDK.EXPECT().Patch(interVrfGwID, interVrfID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(interVrfAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())

		err := resourceNsxtPolicyTier0InterVRFRoutingCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, interVrfID, d.Id())
		assert.Equal(t, interVrfDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(interVrfAPIResponse(), nil)

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())

		err := resourceNsxtPolicyTier0InterVRFRoutingCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails for non-Tier0 gateway path", func(t *testing.T) {
		t1Data := minimalInterVrfData()
		t1Data["gateway_path"] = "/infra/tier-1s/t1-gw-1"

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, t1Data)

		err := resourceNsxtPolicyTier0InterVRFRoutingCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tier0")
	})
}

func TestMockResourceNsxtPolicyTier0InterVRFRoutingRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupInterVrfMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(interVrfAPIResponse(), nil)

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, interVrfDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(nsxModel.PolicyInterVrfRoutingConfig{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyTier0InterVRFRoutingUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupInterVrfMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(interVrfGwID, interVrfID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(interVrfAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0InterVRFRoutingDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupInterVrfMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(interVrfGwID, interVrfID).Return(nil)

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
