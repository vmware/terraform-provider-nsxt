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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	localeservices "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	lsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	evpnTEID          = "evpn-te-1"
	evpnTEName        = "evpn-te-fooname"
	evpnTERevision    = int64(1)
	evpnTEPath        = "/infra/tier-0s/tier0-gw-1/locale-services/ls-1/evpn-tunnel-endpoints/evpn-te-1"
	evpnInterfacePath = "/infra/tier-0s/tier0-gw-1/locale-services/ls-1/interfaces/if-1"
	evpnTEGWID        = "tier0-gw-1"
	evpnTELSID        = "ls-1"
	evpnEdgeNodePath  = "/infra/sites/default/enforcement-points/default/edge-transport-nodes/etn-1/edge-node"
	evpnLocalAddress  = "10.0.0.1"
)

func TestMockResourceNsxtPolicyEvpnTunnelEndpointCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnTESDK := lsmocks.NewMockEvpnTunnelEndpointsClient(ctrl)
	evpnTEWrapper := &localeservices.EvpnTunnelEndpointConfigClientContext{
		Client:     mockEvpnTESDK,
		ClientType: utl.Local,
	}

	originalEvpnTE := cliEvpnTunnelEndpointsClient
	defer func() { cliEvpnTunnelEndpointsClient = originalEvpnTE }()
	cliEvpnTunnelEndpointsClient = func(sessionContext utl.SessionContext, connector client.Connector) *localeservices.EvpnTunnelEndpointConfigClientContext {
		return evpnTEWrapper
	}

	res := resourceNsxtPolicyEvpnTunnelEndpoint()

	t.Run("Create_success", func(t *testing.T) {
		localAddr := evpnLocalAddress
		mockEvpnTESDK.EXPECT().Patch(evpnTEGWID, evpnTELSID, gomock.Any(), gomock.Any()).Return(nil)
		mockEvpnTESDK.EXPECT().Get(evpnTEGWID, evpnTELSID, gomock.Any()).Return(model.EvpnTunnelEndpointConfig{
			DisplayName:    &evpnTEName,
			Path:           &evpnTEPath,
			Revision:       &evpnTERevision,
			LocalAddresses: []string{localAddr},
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":            evpnTEName,
			"external_interface_path": evpnInterfacePath,
			"edge_node_path":          evpnEdgeNodePath,
			"local_address":           evpnLocalAddress,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTunnelEndpointCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, evpnTEName, d.Get("display_name"))
		assert.Equal(t, evpnTEGWID, d.Get("gateway_id"))
		assert.Equal(t, evpnTELSID, d.Get("locale_service_id"))
	})

	t.Run("Create_fails_when_nsx_id_already_exists", func(t *testing.T) {
		mockEvpnTESDK.EXPECT().Get(evpnTEGWID, evpnTELSID, evpnTEID).Return(model.EvpnTunnelEndpointConfig{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":                  evpnTEID,
			"external_interface_path": evpnInterfacePath,
			"edge_node_path":          evpnEdgeNodePath,
			"local_address":           evpnLocalAddress,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTunnelEndpointCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_interface_path_not_T0", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_interface_path": "/infra/tier-1s/gw-id/locale-services/ls-id/interfaces/if-id",
			"edge_node_path":          evpnEdgeNodePath,
			"local_address":           evpnLocalAddress,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTunnelEndpointCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0 External Interface path expected")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockEvpnTESDK.EXPECT().Patch(evpnTEGWID, evpnTELSID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_interface_path": evpnInterfacePath,
			"edge_node_path":          evpnEdgeNodePath,
			"local_address":           evpnLocalAddress,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTunnelEndpointCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEvpnTunnelEndpointRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnTESDK := lsmocks.NewMockEvpnTunnelEndpointsClient(ctrl)
	evpnTEWrapper := &localeservices.EvpnTunnelEndpointConfigClientContext{
		Client:     mockEvpnTESDK,
		ClientType: utl.Local,
	}

	originalEvpnTE := cliEvpnTunnelEndpointsClient
	defer func() { cliEvpnTunnelEndpointsClient = originalEvpnTE }()
	cliEvpnTunnelEndpointsClient = func(sessionContext utl.SessionContext, connector client.Connector) *localeservices.EvpnTunnelEndpointConfigClientContext {
		return evpnTEWrapper
	}

	res := resourceNsxtPolicyEvpnTunnelEndpoint()

	t.Run("Read_success", func(t *testing.T) {
		localAddr := evpnLocalAddress
		mockEvpnTESDK.EXPECT().Get(evpnTEGWID, evpnTELSID, evpnTEID).Return(model.EvpnTunnelEndpointConfig{
			DisplayName:    &evpnTEName,
			Path:           &evpnTEPath,
			Revision:       &evpnTERevision,
			LocalAddresses: []string{localAddr},
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_interface_path": evpnInterfacePath,
		})
		d.Set("gateway_id", evpnTEGWID)
		d.Set("locale_service_id", evpnTELSID)
		d.SetId(evpnTEID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTunnelEndpointRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, evpnTEName, d.Get("display_name"))
		assert.Equal(t, evpnLocalAddress, d.Get("local_address"))
	})

	t.Run("Read_clears_id_when_not_found", func(t *testing.T) {
		mockEvpnTESDK.EXPECT().Get(evpnTEGWID, evpnTELSID, evpnTEID).Return(model.EvpnTunnelEndpointConfig{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_interface_path": evpnInterfacePath,
		})
		d.Set("gateway_id", evpnTEGWID)
		d.Set("locale_service_id", evpnTELSID)
		d.SetId(evpnTEID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTunnelEndpointRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyEvpnTunnelEndpointUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnTESDK := lsmocks.NewMockEvpnTunnelEndpointsClient(ctrl)
	evpnTEWrapper := &localeservices.EvpnTunnelEndpointConfigClientContext{
		Client:     mockEvpnTESDK,
		ClientType: utl.Local,
	}

	originalEvpnTE := cliEvpnTunnelEndpointsClient
	defer func() { cliEvpnTunnelEndpointsClient = originalEvpnTE }()
	cliEvpnTunnelEndpointsClient = func(sessionContext utl.SessionContext, connector client.Connector) *localeservices.EvpnTunnelEndpointConfigClientContext {
		return evpnTEWrapper
	}

	res := resourceNsxtPolicyEvpnTunnelEndpoint()

	t.Run("Update_success", func(t *testing.T) {
		localAddr := evpnLocalAddress
		mockEvpnTESDK.EXPECT().Patch(evpnTEGWID, evpnTELSID, evpnTEID, gomock.Any()).Return(nil)
		mockEvpnTESDK.EXPECT().Get(evpnTEGWID, evpnTELSID, evpnTEID).Return(model.EvpnTunnelEndpointConfig{
			DisplayName:    &evpnTEName,
			Path:           &evpnTEPath,
			Revision:       &evpnTERevision,
			LocalAddresses: []string{localAddr},
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":            evpnTEName,
			"external_interface_path": evpnInterfacePath,
			"edge_node_path":          evpnEdgeNodePath,
			"local_address":           evpnLocalAddress,
		})
		d.Set("gateway_id", evpnTEGWID)
		d.Set("locale_service_id", evpnTELSID)
		d.SetId(evpnTEID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTunnelEndpointUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_Patch_returns_error", func(t *testing.T) {
		mockEvpnTESDK.EXPECT().Patch(evpnTEGWID, evpnTELSID, evpnTEID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_interface_path": evpnInterfacePath,
			"edge_node_path":          evpnEdgeNodePath,
			"local_address":           evpnLocalAddress,
		})
		d.Set("gateway_id", evpnTEGWID)
		d.Set("locale_service_id", evpnTELSID)
		d.SetId(evpnTEID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTunnelEndpointUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEvpnTunnelEndpointDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnTESDK := lsmocks.NewMockEvpnTunnelEndpointsClient(ctrl)
	evpnTEWrapper := &localeservices.EvpnTunnelEndpointConfigClientContext{
		Client:     mockEvpnTESDK,
		ClientType: utl.Local,
	}

	originalEvpnTE := cliEvpnTunnelEndpointsClient
	defer func() { cliEvpnTunnelEndpointsClient = originalEvpnTE }()
	cliEvpnTunnelEndpointsClient = func(sessionContext utl.SessionContext, connector client.Connector) *localeservices.EvpnTunnelEndpointConfigClientContext {
		return evpnTEWrapper
	}

	res := resourceNsxtPolicyEvpnTunnelEndpoint()

	t.Run("Delete_success", func(t *testing.T) {
		mockEvpnTESDK.EXPECT().Delete(evpnTEGWID, evpnTELSID, evpnTEID).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_interface_path": evpnInterfacePath,
		})
		d.Set("gateway_id", evpnTEGWID)
		d.Set("locale_service_id", evpnTELSID)
		d.SetId(evpnTEID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTunnelEndpointDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockEvpnTESDK.EXPECT().Delete(evpnTEGWID, evpnTELSID, evpnTEID).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_interface_path": evpnInterfacePath,
		})
		d.Set("gateway_id", evpnTEGWID)
		d.Set("locale_service_id", evpnTELSID)
		d.SetId(evpnTEID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTunnelEndpointDelete(d, m)
		require.Error(t, err)
	})
}
