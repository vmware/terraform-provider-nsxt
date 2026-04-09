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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	ospfapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services/ospf"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ospfMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services/ospf"
)

var (
	ospfAreaID          = "ospf-area-1"
	ospfAreaDisplayName = "Test OSPF Area"
	ospfAreaDescription = "Test ospf area"
	ospfAreaRevision    = int64(1)
	ospfAreaOspfPath    = "/infra/tier-0s/t0-gw-1/locale-services/ls-1/ospf"
	ospfAreaGwID        = "t0-gw-1"
	ospfAreaLsID        = "ls-1"
	ospfAreaID2         = "10"
)

func ospfAreaAPIResponse() nsxModel.OspfAreaConfig {
	areaType := nsxModel.OspfAreaConfig_AREA_TYPE_NSSA
	return nsxModel.OspfAreaConfig{
		Id:          &ospfAreaID,
		DisplayName: &ospfAreaDisplayName,
		Description: &ospfAreaDescription,
		Revision:    &ospfAreaRevision,
		AreaType:    &areaType,
		AreaId:      &ospfAreaID2,
	}
}

func minimalOSPFAreaData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": ospfAreaDisplayName,
		"description":  ospfAreaDescription,
		"nsx_id":       ospfAreaID,
		"ospf_path":    ospfAreaOspfPath,
		"area_id":      ospfAreaID2,
		"area_type":    nsxModel.OspfAreaConfig_AREA_TYPE_NSSA,
		"auth_mode":    nsxModel.OspfAuthenticationConfig_MODE_NONE,
	}
}

func setupOSPFAreaMock(t *testing.T, ctrl *gomock.Controller) (*ospfMocks.MockAreasClient, func()) {
	mockSDK := ospfMocks.NewMockAreasClient(ctrl)
	mockWrapper := &ospfapi.OspfAreaConfigClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliOspfAreasClient
	cliOspfAreasClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *ospfapi.OspfAreaConfigClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliOspfAreasClient = original }
}

func TestMockResourceNsxtPolicyOspfAreaCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupOSPFAreaMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(ospfAreaGwID, ospfAreaLsID, ospfAreaID, gomock.Any()).Return(ospfAreaAPIResponse(), nil),
			mockSDK.EXPECT().Get(ospfAreaGwID, ospfAreaLsID, ospfAreaID).Return(ospfAreaAPIResponse(), nil),
		)

		res := resourceNsxtPolicyOspfArea()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFAreaData())

		err := resourceNsxtPolicyOspfAreaCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ospfAreaID, d.Id())
	})
}

func TestMockResourceNsxtPolicyOspfAreaRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupOSPFAreaMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ospfAreaGwID, ospfAreaLsID, ospfAreaID).Return(ospfAreaAPIResponse(), nil)

		res := resourceNsxtPolicyOspfArea()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFAreaData())
		d.SetId(ospfAreaID)

		err := resourceNsxtPolicyOspfAreaRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ospfAreaDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(ospfAreaGwID, ospfAreaLsID, ospfAreaID).Return(nsxModel.OspfAreaConfig{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyOspfArea()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFAreaData())
		d.SetId(ospfAreaID)

		err := resourceNsxtPolicyOspfAreaRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyOspfArea()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFAreaData())

		err := resourceNsxtPolicyOspfAreaRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyOspfAreaUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupOSPFAreaMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(ospfAreaGwID, ospfAreaLsID, ospfAreaID, gomock.Any()).Return(ospfAreaAPIResponse(), nil),
			mockSDK.EXPECT().Get(ospfAreaGwID, ospfAreaLsID, ospfAreaID).Return(ospfAreaAPIResponse(), nil),
		)

		res := resourceNsxtPolicyOspfArea()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFAreaData())
		d.SetId(ospfAreaID)

		err := resourceNsxtPolicyOspfAreaUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyOspfArea()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFAreaData())

		err := resourceNsxtPolicyOspfAreaUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyOspfAreaDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupOSPFAreaMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ospfAreaGwID, ospfAreaLsID, ospfAreaID).Return(nil)

		res := resourceNsxtPolicyOspfArea()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFAreaData())
		d.SetId(ospfAreaID)

		err := resourceNsxtPolicyOspfAreaDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyOspfArea()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFAreaData())

		err := resourceNsxtPolicyOspfAreaDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
