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

	transportnodecollections "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points/transport_node_collections"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tncmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points/transport_node_collections"
)

var tncRealizationPath = "/infra/sites/default/enforcement-points/default/transport-node-collections/tnc-1"

func setupTransportNodeCollectionStateMock(t *testing.T, ctrl *gomock.Controller) (*tncmocks.MockStateClient, func()) {
	mockSDK := tncmocks.NewMockStateClient(ctrl)
	mockWrapper := &transportnodecollections.StateClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliTransportNodeCollectionStateClient
	cliTransportNodeCollectionStateClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *transportnodecollections.StateClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliTransportNodeCollectionStateClient = original }
}

func TestMockDataSourceNsxtPolicyHostTransportNodeCollectionRealizationRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTransportNodeCollectionStateMock(t, ctrl)
	defer restore()

	ds := dataSourceNsxtPolicyHostTransportNodeCollectionRealization()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get("default", "default", "tnc-1").Return(nsxModel.TransportNodeCollectionState{
			State: str(nsxModel.TransportNodeCollectionState_STATE_SUCCESS),
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":  tncRealizationPath,
			"delay": 0,
		})

		err := dataSourceNsxtPolicyHostTransportNodeCollectionRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "tnc-1", d.Id())
		assert.Equal(t, nsxModel.TransportNodeCollectionState_STATE_SUCCESS, d.Get("state"))
	})

	t.Run("Read failure state", func(t *testing.T) {
		mockSDK.EXPECT().Get("default", "default", "tnc-1").Return(nsxModel.TransportNodeCollectionState{
			State: str(nsxModel.TransportNodeCollectionState_STATE_FAILED_TO_REALIZE),
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    tncRealizationPath,
			"delay":   0,
			"timeout": 5,
		})

		err := dataSourceNsxtPolicyHostTransportNodeCollectionRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read API error", func(t *testing.T) {
		mockSDK.EXPECT().Get("default", "default", "tnc-1").Return(nsxModel.TransportNodeCollectionState{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    tncRealizationPath,
			"delay":   0,
			"timeout": 5,
		})

		err := dataSourceNsxtPolicyHostTransportNodeCollectionRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read invalid path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": "/infra/invalid-path",
		})

		err := dataSourceNsxtPolicyHostTransportNodeCollectionRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid transport node collection path")
	})
}
