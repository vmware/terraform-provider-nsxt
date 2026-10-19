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
// The lifecycle test below also covers success with a serialized HostSwitchSpec.
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
// The lifecycle test below covers the success path.
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
// The lifecycle test below covers Update followed by Read.
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

func TestMockResourceNsxtPolicyHostTransportNodeProfileLifecycleSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	sdk, restore := setupHtnpMock(t, ctrl)
	t.Cleanup(restore)
	m := newIsolatedMockProviderClient(t)
	res := resourceNsxtPolicyHostTransportNodeProfile()
	d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnpData())
	spec, err := getHostSwitchSpecFromSchema(d, m, nodeTypeHost)
	require.NoError(t, err)
	obj := nsxModel.PolicyHostTransportNodeProfile{Id: &htnpID, DisplayName: &htnpDisplayName, Description: &htnpDescription, HostSwitchSpec: spec, IgnoreOverriddenHosts: boolPtr(true)}
	gomock.InOrder(
		sdk.EXPECT().Get(htnpID).Return(nsxModel.PolicyHostTransportNodeProfile{}, vapiErrors.NotFound{}),
		sdk.EXPECT().Update(htnpID, gomock.Any(), nil, nil).DoAndReturn(func(_ string, got nsxModel.PolicyHostTransportNodeProfile, _ *bool, _ *bool) (nsxModel.PolicyHostTransportNodeProfile, error) {
			require.Equal(t, htnpDisplayName, *got.DisplayName)
			require.NotNil(t, got.HostSwitchSpec)
			return obj, nil
		}),
		sdk.EXPECT().Get(htnpID).Return(obj, nil),
	)
	require.NoError(t, resourceNsxtPolicyHostTransportNodeProfileCreate(d, m))
	assert.Equal(t, htnpID, d.Id())
	assert.Equal(t, htnpDisplayName, d.Get("display_name"))
	sdk.EXPECT().Update(htnpID, gomock.Any(), nil, nil).Return(obj, nil)
	sdk.EXPECT().Get(htnpID).Return(obj, nil)
	require.NoError(t, resourceNsxtPolicyHostTransportNodeProfileUpdate(d, m))
}
