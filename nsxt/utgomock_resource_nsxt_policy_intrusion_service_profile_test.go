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

	intrusionservices "github.com/vmware/terraform-provider-nsxt/api/infra/settings/firewall/security/intrusion_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	isprofmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/settings/firewall/security/intrusion_services"
)

var (
	isprofID          = "ids-profile-1"
	isprofDisplayName = "Test IDS Profile"
	isprofDescription = "test intrusion service profile"
	isprofRevision    = int64(1)
	isprofPath        = "/infra/settings/firewall/security/intrusion-services/profiles/ids-profile-1"
	isprofSeverity    = nsxModel.IdsProfile_PROFILE_SEVERITY_HIGH
)

func isprofAPIResponse() nsxModel.IdsProfile {
	return nsxModel.IdsProfile{
		Id:              &isprofID,
		DisplayName:     &isprofDisplayName,
		Description:     &isprofDescription,
		Revision:        &isprofRevision,
		Path:            &isprofPath,
		ProfileSeverity: []string{isprofSeverity},
	}
}

func minimalIsprofData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": isprofDisplayName,
		"description":  isprofDescription,
		"nsx_id":       isprofID,
		"severities":   []interface{}{isprofSeverity},
	}
}

func setupIsprofMock(t *testing.T, ctrl *gomock.Controller) (*isprofmocks.MockProfilesClient, func()) {
	mockSDK := isprofmocks.NewMockProfilesClient(ctrl)
	mockWrapper := &intrusionservices.IdsProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliIntrusionServiceProfilesClient
	cliIntrusionServiceProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *intrusionservices.IdsProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliIntrusionServiceProfilesClient = original }
}

func TestMockResourceNsxtPolicyIntrusionServiceProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIsprofMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(isprofID).Return(nsxModel.IdsProfile{}, notFoundErr),
			mockSDK.EXPECT().Patch(isprofID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(isprofID).Return(isprofAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())

		err := resourceNsxtPolicyIntrusionServiceProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, isprofID, d.Id())
		assert.Equal(t, isprofDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(isprofID).Return(isprofAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())

		err := resourceNsxtPolicyIntrusionServiceProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIsprofMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(isprofID).Return(isprofAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())
		d.SetId(isprofID)

		err := resourceNsxtPolicyIntrusionServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, isprofDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(isprofID).Return(nsxModel.IdsProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())
		d.SetId(isprofID)

		err := resourceNsxtPolicyIntrusionServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())

		err := resourceNsxtPolicyIntrusionServiceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIsprofMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(isprofID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(isprofID).Return(isprofAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())
		d.SetId(isprofID)

		err := resourceNsxtPolicyIntrusionServiceProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())

		err := resourceNsxtPolicyIntrusionServiceProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIsprofMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(isprofID).Return(nil)

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())
		d.SetId(isprofID)

		err := resourceNsxtPolicyIntrusionServiceProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())

		err := resourceNsxtPolicyIntrusionServiceProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
