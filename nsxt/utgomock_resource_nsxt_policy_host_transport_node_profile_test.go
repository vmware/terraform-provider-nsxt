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

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	htnpID          = "htnp-1"
	htnpDisplayName = "Test HTN Profile"
	htnpDescription = "test host transport node profile"
)

func minimalHtnpData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": htnpDisplayName,
		"description":  htnpDescription,
		"nsx_id":       htnpID,
	}
}

func setupHtnpMock(t *testing.T, ctrl *gomock.Controller) (*inframocks.MockHostTransportNodeProfilesClient, func()) {
	mockSDK := inframocks.NewMockHostTransportNodeProfilesClient(ctrl)
	mockWrapper := &infraapi.PolicyHostTransportNodeProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliHostTransportNodeProfilesClient
	cliHostTransportNodeProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.PolicyHostTransportNodeProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliHostTransportNodeProfilesClient = original }
}

// TestMockResourceNsxtPolicyHostTransportNodeProfileCreate tests the Create function.
// The full success path (including Read) is not testable because setHostSwitchSpecInSchema
// requires a non-nil, properly serialized HostSwitchSpec StructValue.
func TestMockResourceNsxtPolicyHostTransportNodeProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtnpMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		existingProfile := nsxModel.PolicyHostTransportNodeProfile{
			Id:          &htnpID,
			DisplayName: &htnpDisplayName,
		}
		mockSDK.EXPECT().Get(htnpID).Return(existingProfile, nil)

		res := resourceNsxtPolicyHostTransportNodeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnpData())

		err := resourceNsxtPolicyHostTransportNodeProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

// TestMockResourceNsxtPolicyHostTransportNodeProfileRead tests the Read function.
// The success path is not testable because setHostSwitchSpecInSchema requires
// a non-nil, properly serialized HostSwitchSpec StructValue that is complex
// to construct in unit tests.
func TestMockResourceNsxtPolicyHostTransportNodeProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtnpMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(htnpID).Return(nsxModel.PolicyHostTransportNodeProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyHostTransportNodeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnpData())
		d.SetId(htnpID)

		err := resourceNsxtPolicyHostTransportNodeProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyHostTransportNodeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnpData())

		err := resourceNsxtPolicyHostTransportNodeProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(htnpID).Return(nsxModel.PolicyHostTransportNodeProfile{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyHostTransportNodeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnpData())
		d.SetId(htnpID)

		err := resourceNsxtPolicyHostTransportNodeProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

// TestMockResourceNsxtPolicyHostTransportNodeProfileUpdate tests the Update function.
// The success path calls Read after Update, which is not fully testable due to
// the same HostSwitchSpec serialization constraint.
func TestMockResourceNsxtPolicyHostTransportNodeProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtnpMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyHostTransportNodeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnpData())

		err := resourceNsxtPolicyHostTransportNodeProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Update(htnpID, gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.PolicyHostTransportNodeProfile{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyHostTransportNodeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnpData())
		d.SetId(htnpID)

		err := resourceNsxtPolicyHostTransportNodeProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyHostTransportNodeProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtnpMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(htnpID).Return(nil)

		res := resourceNsxtPolicyHostTransportNodeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnpData())
		d.SetId(htnpID)

		err := resourceNsxtPolicyHostTransportNodeProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyHostTransportNodeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnpData())

		err := resourceNsxtPolicyHostTransportNodeProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
