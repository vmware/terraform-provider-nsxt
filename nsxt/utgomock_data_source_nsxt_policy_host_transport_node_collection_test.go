//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Reuses setupHtncMock/htncAPIResponse/htnc* fixtures defined in
// utgomock_resource_nsxt_policy_host_transport_node_collection_test.go, since
// both the resource and this data source use cliTransportNodeCollectionsClient.

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"
)

func TestMockDataSourceNsxtPolicyHostTransportNodeCollectionRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtncMock(t, ctrl)
	defer restore()

	ds := dataSourceNsxtPolicyHostTransportNodeCollection()

	newClient := func() nsxtClients {
		c := newGoMockProviderClient()
		c.PolicyEnforcementPoint = htncEPID
		return c
	}

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(defaultSite, htncEPID, htncID).Return(htncAPIResponse(), nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": htncID,
		})

		err := dataSourceNsxtPolicyHostTransportNodeCollectionRead(d, newClient())
		require.NoError(t, err)
		assert.Equal(t, htncID, d.Id())
		assert.Equal(t, htncDisplayName, d.Get("display_name"))
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(defaultSite, htncEPID, htncID).Return(nsxModel.HostTransportNodeCollection{}, errors.New("get failed"))

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": htncID,
		})

		err := dataSourceNsxtPolicyHostTransportNodeCollectionRead(d, newClient())
		require.Error(t, err)
	})

	t.Run("by display_name match", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, htncEPID, nil, nil, nil).Return(nsxModel.HostTransportNodeCollectionListResult{
			Results: []nsxModel.HostTransportNodeCollection{htncAPIResponse()},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": htncDisplayName,
		})

		err := dataSourceNsxtPolicyHostTransportNodeCollectionRead(d, newClient())
		require.NoError(t, err)
		assert.Equal(t, htncID, d.Id())
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, htncEPID, nil, nil, nil).Return(nsxModel.HostTransportNodeCollectionListResult{}, errors.New("list failed"))

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": htncDisplayName,
		})

		err := dataSourceNsxtPolicyHostTransportNodeCollectionRead(d, newClient())
		require.Error(t, err)
	})

	t.Run("no match from list", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, htncEPID, nil, nil, nil).Return(nsxModel.HostTransportNodeCollectionListResult{
			Results: []nsxModel.HostTransportNodeCollection{},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtPolicyHostTransportNodeCollectionRead(d, newClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("global manager not supported", func(t *testing.T) {
		c := newClient()
		c.PolicyGlobalManager = true

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": htncID,
		})

		err := dataSourceNsxtPolicyHostTransportNodeCollectionRead(d, c)
		require.Error(t, err)
	})
}
