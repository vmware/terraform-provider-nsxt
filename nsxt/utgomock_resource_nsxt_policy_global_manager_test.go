//go:build unittest

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
	gmModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	gmMocks "github.com/vmware/terraform-provider-nsxt/mocks/global_infra"
)

var (
	gmID          = "gm-1"
	gmDisplayName = "Test Global Manager"
	gmDescription = "test global manager"
	gmRevision    = int64(1)
	gmPath        = "/global-infra/global-managers/gm-1"
	gmMode        = gmModel.GlobalManager_MODE_ACTIVE
)

func globalManagerAPIResponse() gmModel.GlobalManager {
	failIfRtt := true
	maxRtt := int64(250)
	return gmModel.GlobalManager{
		Id:                &gmID,
		DisplayName:       &gmDisplayName,
		Description:       &gmDescription,
		Revision:          &gmRevision,
		Path:              &gmPath,
		Mode:              &gmMode,
		FailIfRttExceeded: &failIfRtt,
		MaximumRtt:        &maxRtt,
	}
}

func minimalGlobalManagerData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": gmDisplayName,
		"description":  gmDescription,
		"nsx_id":       gmID,
		"mode":         gmMode,
	}
}

func setupGlobalManagerMock(t *testing.T, ctrl *gomock.Controller) (*gmMocks.MockGlobalManagersClient, func()) {
	mockSDK := gmMocks.NewMockGlobalManagersClient(ctrl)
	mockWrapper := &infraapi.GlobalManagerClientContext{
		Client:     mockSDK,
		ClientType: utl.Global,
	}

	original := cliGlobalManagersClient
	cliGlobalManagersClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.GlobalManagerClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliGlobalManagersClient = original }
}

func TestMockResourceNsxtPolicyGlobalManagerGuard(t *testing.T) {
	t.Run("Create fails when not global manager", func(t *testing.T) {
		res := resourceNsxtPolicyGlobalManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGlobalManagerData())

		err := resourceNsxtPolicyGlobalManagerCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})
}

func TestMockResourceNsxtPolicyGlobalManagerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGlobalManagerMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gmID).Return(gmModel.GlobalManager{}, notFoundErr),
			mockSDK.EXPECT().Patch(gmID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gmID).Return(globalManagerAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGlobalManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGlobalManagerData())

		err := resourceNsxtPolicyGlobalManagerCreate(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
		assert.Equal(t, gmID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(gmID).Return(globalManagerAPIResponse(), nil)

		res := resourceNsxtPolicyGlobalManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGlobalManagerData())

		err := resourceNsxtPolicyGlobalManagerCreate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGlobalManagerRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGlobalManagerMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(gmID).Return(globalManagerAPIResponse(), nil)

		res := resourceNsxtPolicyGlobalManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGlobalManagerData())
		d.SetId(gmID)

		err := resourceNsxtPolicyGlobalManagerRead(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
		assert.Equal(t, gmDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(gmID).Return(gmModel.GlobalManager{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGlobalManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGlobalManagerData())
		d.SetId(gmID)

		err := resourceNsxtPolicyGlobalManagerRead(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGlobalManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGlobalManagerData())

		err := resourceNsxtPolicyGlobalManagerRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGlobalManagerUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGlobalManagerMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Update(gmID, gomock.Any(), gomock.Any()).Return(globalManagerAPIResponse(), nil)

		res := resourceNsxtPolicyGlobalManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGlobalManagerData())
		d.SetId(gmID)

		err := resourceNsxtPolicyGlobalManagerUpdate(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGlobalManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGlobalManagerData())

		err := resourceNsxtPolicyGlobalManagerUpdate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGlobalManagerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGlobalManagerMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(gmID).Return(nil)

		res := resourceNsxtPolicyGlobalManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGlobalManagerData())
		d.SetId(gmID)

		err := resourceNsxtPolicyGlobalManagerDelete(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGlobalManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGlobalManagerData())

		err := resourceNsxtPolicyGlobalManagerDelete(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
	})
}
