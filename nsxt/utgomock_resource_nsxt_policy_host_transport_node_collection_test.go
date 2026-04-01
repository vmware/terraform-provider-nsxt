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

	enforcement_points "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
)

var (
	htncID          = "htnc-1"
	htncDisplayName = "Test HTNC"
	htncDescription = "test host transport node collection"
	htncRevision    = int64(1)
	htncSitePath    = "/infra/sites/default"
	htncSiteID      = "default"
	htncEPID        = "default"
	htncPath        = "/infra/sites/default/enforcement-points/default/transport-node-collections/htnc-1"
	htncCCID        = "cc-domain-c1"
)

func htncAPIResponse() nsxModel.HostTransportNodeCollection {
	return nsxModel.HostTransportNodeCollection{
		Id:                  &htncID,
		DisplayName:         &htncDisplayName,
		Description:         &htncDescription,
		Revision:            &htncRevision,
		Path:                &htncPath,
		ComputeCollectionId: &htncCCID,
	}
}

func minimalHtncData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":          htncDisplayName,
		"description":           htncDescription,
		"nsx_id":                htncID,
		"site_path":             htncSitePath,
		"enforcement_point":     htncEPID,
		"compute_collection_id": htncCCID,
	}
}

func setupHtncMock(t *testing.T, ctrl *gomock.Controller) (*epmocks.MockTransportNodeCollectionsClient, func()) {
	mockSDK := epmocks.NewMockTransportNodeCollectionsClient(ctrl)
	mockWrapper := &enforcement_points.HostTransportNodeCollectionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliTransportNodeCollectionsClient
	cliTransportNodeCollectionsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *enforcement_points.HostTransportNodeCollectionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTransportNodeCollectionsClient = original }
}

func TestMockResourceNsxtPolicyHostTransportNodeCollectionRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtncMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(htncSiteID, htncEPID, htncID).Return(htncAPIResponse(), nil)

		res := resourceNsxtPolicyHostTransportNodeCollection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtncData())
		d.SetId(htncID)

		err := resourceNsxtPolicyHostTransportNodeCollectionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, htncDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(htncSiteID, htncEPID, htncID).Return(nsxModel.HostTransportNodeCollection{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyHostTransportNodeCollection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtncData())
		d.SetId(htncID)

		err := resourceNsxtPolicyHostTransportNodeCollectionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})
}

func TestMockResourceNsxtPolicyHostTransportNodeCollectionCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHtncSDK, restoreHtnc := setupHtncMock(t, ctrl)
	defer restoreHtnc()
	// Site check is needed during exists check
	mockSiteSDK, restoreSite := setupHtnSiteMock(t, ctrl)
	defer restoreSite()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSiteSDK.EXPECT().Get(htncSiteID).Return(nsxModel.Site{}, nil),
			mockHtncSDK.EXPECT().Get(htncSiteID, htncEPID, htncID).Return(nsxModel.HostTransportNodeCollection{}, vapiErrors.NotFound{}),
			mockHtncSDK.EXPECT().Update(htncSiteID, htncEPID, htncID, gomock.Any(), gomock.Any(), gomock.Any()).Return(htncAPIResponse(), nil),
			mockHtncSDK.EXPECT().Get(htncSiteID, htncEPID, htncID).Return(htncAPIResponse(), nil),
		)

		res := resourceNsxtPolicyHostTransportNodeCollection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtncData())

		err := resourceNsxtPolicyHostTransportNodeCollectionCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, htncID, d.Id())
	})
}

func TestMockResourceNsxtPolicyHostTransportNodeCollectionUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtncMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(htncSiteID, htncEPID, htncID, gomock.Any(), gomock.Any(), gomock.Any()).Return(htncAPIResponse(), nil),
			mockSDK.EXPECT().Get(htncSiteID, htncEPID, htncID).Return(htncAPIResponse(), nil),
		)

		res := resourceNsxtPolicyHostTransportNodeCollection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtncData())
		d.SetId(htncID)

		err := resourceNsxtPolicyHostTransportNodeCollectionUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyHostTransportNodeCollectionDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtncMock(t, ctrl)
	defer restore()

	t.Run("Delete success (without NSX removal)", func(t *testing.T) {
		mockSDK.EXPECT().Delete(htncSiteID, htncEPID, htncID).Return(nil)

		res := resourceNsxtPolicyHostTransportNodeCollection()
		data := minimalHtncData()
		data["remove_nsx_on_destroy"] = false
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(htncID)

		err := resourceNsxtPolicyHostTransportNodeCollectionDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
