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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"

	fabricapi "github.com/vmware/terraform-provider-nsxt/api/nsx/fabric"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	fabricmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/fabric"
)

var (
	ccTestExternalID  = "cc-ext-id-1"
	ccTestName        = "test-compute-collection"
	ccTestOriginType  = "VC_Cluster"
	ccTestOriginID    = "vcenter-id-1"
	ccTestCmLocalID   = "domain-123"
	ccTestResultCount = int64(1)
)

func computeCollectionModel() nsxModel.ComputeCollection {
	return nsxModel.ComputeCollection{
		ExternalId:  &ccTestExternalID,
		DisplayName: &ccTestName,
		OriginType:  &ccTestOriginType,
		OriginId:    &ccTestOriginID,
		CmLocalId:   &ccTestCmLocalID,
	}
}

func setupComputeCollectionMock(t *testing.T, ctrl *gomock.Controller) (*fabricmocks.MockComputeCollectionsClient, func()) {
	t.Helper()
	mockSDK := fabricmocks.NewMockComputeCollectionsClient(ctrl)
	wrapper := &fabricapi.ComputeCollectionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliComputeCollectionsClient
	cliComputeCollectionsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *fabricapi.ComputeCollectionClientContext {
		return wrapper
	}
	return mockSDK, func() { cliComputeCollectionsClient = orig }
}

func TestMockDataSourceNsxtComputeCollectionRead(t *testing.T) {
	t.Run("by ID success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeCollectionMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(ccTestExternalID).Return(computeCollectionModel(), nil)

		ds := dataSourceNsxtComputeCollection()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": ccTestExternalID,
		})

		err := dataSourceNsxtComputeCollectionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ccTestExternalID, d.Id())
		assert.Equal(t, ccTestName, d.Get("display_name"))
		assert.Equal(t, ccTestOriginType, d.Get("origin_type"))
	})

	t.Run("by ID API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeCollectionMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(ccTestExternalID).Return(nsxModel.ComputeCollection{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtComputeCollection()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": ccTestExternalID,
		})

		err := dataSourceNsxtComputeCollectionRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read ComputeCollection")
	})

	t.Run("by name single match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeCollectionMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().List(nil, nil, nil, &ccTestName, nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ComputeCollectionListResult{
			ResultCount: &ccTestResultCount,
			Results:     []nsxModel.ComputeCollection{computeCollectionModel()},
		}, nil)

		ds := dataSourceNsxtComputeCollection()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": ccTestName,
		})

		err := dataSourceNsxtComputeCollectionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ccTestExternalID, d.Id())
	})

	t.Run("by name list error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeCollectionMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.ComputeCollectionListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtComputeCollection()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": ccTestName,
		})

		err := dataSourceNsxtComputeCollectionRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read Compute Collections")
	})
}
