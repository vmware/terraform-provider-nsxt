//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mpmodel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"
)

// setupStateClientMock and tnID are defined in
// utgomock_resource_nsxt_edge_transport_node_test.go and reused here.

func transportNodeRealizationData(id string) map[string]interface{} {
	return map[string]interface{}{
		"id":      id,
		"timeout": 5,
		"delay":   0,
	}
}

func TestUnitNsxt_DataSourceNsxtTransportNodeRealizationRead(t *testing.T) {
	t.Run("reaches SUCCESS state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupStateClientMock(t, ctrl)
		defer restore()

		successState := mpmodel.TransportNodeState_STATE_SUCCESS
		mockSDK.EXPECT().Get(tnID).Return(mpmodel.TransportNodeState{State: &successState}, nil)

		ds := dataSourceNsxtTransportNodeRealization()
		d := schema.TestResourceDataRaw(t, ds.Schema, transportNodeRealizationData(tnID))

		err := dataSourceNsxtTransportNodeRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, mpmodel.TransportNodeState_STATE_SUCCESS, d.Get("state"))
	})

	t.Run("reaches FAILED state and errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupStateClientMock(t, ctrl)
		defer restore()

		failedState := mpmodel.TransportNodeState_STATE_FAILED
		failureMsg := "realization failed"
		mockSDK.EXPECT().Get(tnID).Return(mpmodel.TransportNodeState{
			State:          &failedState,
			FailureMessage: &failureMsg,
		}, nil).AnyTimes()

		ds := dataSourceNsxtTransportNodeRealization()
		d := schema.TestResourceDataRaw(t, ds.Schema, transportNodeRealizationData(tnID))

		err := dataSourceNsxtTransportNodeRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to realize")
	})

	t.Run("API error while polling is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupStateClientMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tnID).Return(mpmodel.TransportNodeState{}, errors.New("get state failed")).AnyTimes()

		ds := dataSourceNsxtTransportNodeRealization()
		d := schema.TestResourceDataRaw(t, ds.Schema, transportNodeRealizationData(tnID))

		err := dataSourceNsxtTransportNodeRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
