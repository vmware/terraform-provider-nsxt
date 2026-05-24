//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/nsx/TransportNodesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/TransportNodesClient.go TransportNodesClient
// mockgen -destination=mocks/nsx/transport_nodes/StateClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/transport_nodes/StateClient.go StateClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx"
	mpmodel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/transport_nodes"
	"go.uber.org/mock/gomock"

	nsxmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx"
	statemocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/transport_nodes"
)

var (
	tnID          = "transport-node-uuid-1"
	tnDisplayName = "test-transport-node"
	tnDescription = "Test transport node"
	tnRevision    = int64(1)
	tnFQDN        = "edge-node.example.com"
	tnExternalID  = "external-id-1"
)

// edgeNodeStructValue builds the *data.StructValue the resource's Read path
// expects in TransportNode.NodeDeploymentInfo.
// NodeSettings is intentionally omitted: Read no longer uses it because the
// NSX edge MPA initial-boot sync can temporarily overwrite the desired state
// in the API with the edge VM's current running state, causing perpetual drift
// if those stale values are written back to Terraform state (bug 3689817).
func edgeNodeStructValue() *data.StructValue {
	node := mpmodel.EdgeNode{
		ResourceType: mpmodel.EdgeNode__TYPE_IDENTIFIER,
		ExternalId:   &tnExternalID,
		Fqdn:         &tnFQDN,
	}
	converter := bindings.NewTypeConverter()
	dataValue, errs := converter.ConvertToVapi(node, mpmodel.EdgeNodeBindingType())
	if errs != nil {
		panic(errs[0])
	}
	return dataValue.(*data.StructValue)
}

// edgeNodeStructValueWithNodeSettings returns a StructValue whose NodeSettings
// carries a hostname that differs from what the Terraform config would hold.
// Used to verify that Read does NOT copy MPA-sync-overwritten values into state.
func edgeNodeStructValueWithNodeSettings(apiHostname string) *data.StructValue {
	node := mpmodel.EdgeNode{
		ResourceType: mpmodel.EdgeNode__TYPE_IDENTIFIER,
		ExternalId:   &tnExternalID,
		Fqdn:         &tnFQDN,
		NodeSettings: &mpmodel.EdgeNodeSettings{
			Hostname: &apiHostname,
		},
	}
	converter := bindings.NewTypeConverter()
	dataValue, errs := converter.ConvertToVapi(node, mpmodel.EdgeNodeBindingType())
	if errs != nil {
		panic(errs[0])
	}
	return dataValue.(*data.StructValue)
}

func transportNodeAPIResponse() mpmodel.TransportNode {
	return mpmodel.TransportNode{
		Id:                 &tnID,
		DisplayName:        &tnDisplayName,
		Description:        &tnDescription,
		Revision:           &tnRevision,
		NodeDeploymentInfo: edgeNodeStructValue(),
	}
}

func setupTransportNodeMock(t *testing.T, ctrl *gomock.Controller) (*nsxmocks.MockTransportNodesClient, func()) {
	mockSDK := nsxmocks.NewMockTransportNodesClient(ctrl)

	originalCli := cliTransportNodesClient
	cliTransportNodesClient = func(_ client.Connector) nsx.TransportNodesClient {
		return mockSDK
	}
	return mockSDK, func() { cliTransportNodesClient = originalCli }
}

func setupStateClientMock(t *testing.T, ctrl *gomock.Controller) (*statemocks.MockStateClient, func()) {
	mockSDK := statemocks.NewMockStateClient(ctrl)

	originalCli := cliTransportNodeStateClient
	cliTransportNodeStateClient = func(_ client.Connector) transport_nodes.StateClient {
		return mockSDK
	}
	return mockSDK, func() { cliTransportNodeStateClient = originalCli }
}

func minimalTransportNodeData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": tnDisplayName,
		"description":  tnDescription,
		"node_id":      tnID,
	}
}

func transportNodeDataWithNodeSettings() map[string]interface{} {
	m := minimalTransportNodeData()
	m["revision"] = int(tnRevision)
	m["node_settings"] = []interface{}{
		map[string]interface{}{
			"hostname":       tnFQDN,
			"dns_servers":    []interface{}{"10.0.0.1"},
			"search_domains": []interface{}{"example.com"},
		},
	}
	return m
}

func TestMockResourceNsxtEdgeTransportNodeCreate(t *testing.T) {
	t.Run("Create with existing node_id success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		resp := transportNodeAPIResponse()
		// Get in Create, Get from Read after Update (Update does not pre-Get)
		mockSDK.EXPECT().Get(tnID).Return(resp, nil).Times(2)
		mockSDK.EXPECT().Update(tnID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(resp, nil)

		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTransportNodeData())

		err := resourceNsxtEdgeTransportNodeCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tnID, d.Id())
	})

	t.Run("Create with existing node_id fails when Get returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tnID).Return(mpmodel.TransportNode{}, errors.New("get API error"))

		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTransportNodeData())

		err := resourceNsxtEdgeTransportNodeCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "get API error")
	})
}

func TestMockResourceNsxtEdgeTransportNodeRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tnID).Return(transportNodeAPIResponse(), nil)

		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(tnID)

		err := resourceNsxtEdgeTransportNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tnDisplayName, d.Get("display_name"))
		assert.Equal(t, tnDescription, d.Get("description"))
		assert.Equal(t, int(tnRevision), d.Get("revision"))
		assert.Equal(t, tnExternalID, d.Get("external_id"))
		assert.Equal(t, tnFQDN, d.Get("fqdn"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtEdgeTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Read fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tnID).Return(mpmodel.TransportNode{}, errors.New("read API error"))

		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(tnID)

		err := resourceNsxtEdgeTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read API error")
	})

	t.Run("Read does not overwrite node_settings from API", func(t *testing.T) {
		// When a freshly deployed edge VM boots, its MPA performs an initial
		// sync that temporarily overwrites the desired node_settings in the NSX
		// API with the VM's current running state (bug 3689817).  Read must NOT
		// copy those MPA-overwritten values into Terraform state or it would
		// cause perpetual drift.  This test verifies that the config-driven
		// node_settings value in state survives a Read unchanged, even when the
		// API response carries a different (stale) hostname.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		apiHostname := "stale-mpa-reported-hostname.example.com"
		apiResp := mpmodel.TransportNode{
			Id:                 &tnID,
			DisplayName:        &tnDisplayName,
			Description:        &tnDescription,
			Revision:           &tnRevision,
			NodeDeploymentInfo: edgeNodeStructValueWithNodeSettings(apiHostname),
		}
		mockSDK.EXPECT().Get(tnID).Return(apiResp, nil)

		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, transportNodeDataWithNodeSettings())
		d.SetId(tnID)

		err := resourceNsxtEdgeTransportNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		// The config-driven hostname (tnFQDN) must survive Read unchanged.
		// The stale MPA-reported hostname must NOT replace it.
		nodeSettings := d.Get("node_settings").([]interface{})
		require.Len(t, nodeSettings, 1)
		ns := nodeSettings[0].(map[string]interface{})
		assert.Equal(t, tnFQDN, ns["hostname"],
			"Read must not overwrite node_settings.hostname with MPA-reported stale value")
		assert.NotEqual(t, apiHostname, ns["hostname"])
	})
}

func TestMockResourceNsxtEdgeTransportNodeUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		resp := transportNodeAPIResponse()
		mockSDK.EXPECT().Update(tnID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(resp, nil)
		mockSDK.EXPECT().Get(tnID).Return(resp, nil)

		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTransportNodeData())
		d.SetId(tnID)

		err := resourceNsxtEdgeTransportNodeUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tnDisplayName, d.Get("display_name"))
	})

	t.Run("Update with node_settings when node_id is set", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		resp := transportNodeAPIResponse()
		mockSDK.EXPECT().Update(tnID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(resp, nil)
		mockSDK.EXPECT().Get(tnID).Return(resp, nil)

		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, transportNodeDataWithNodeSettings())
		d.SetId(tnID)

		err := resourceNsxtEdgeTransportNodeUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTransportNodeData())

		err := resourceNsxtEdgeTransportNodeUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Update fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(tnID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mpmodel.TransportNode{}, errors.New("update API error"))

		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTransportNodeData())
		d.SetId(tnID)

		err := resourceNsxtEdgeTransportNodeUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update API error")
	})
}

func TestMockResourceNsxtEdgeTransportNodeDelete(t *testing.T) {
	// Override state poll delays so the WaitForState loop completes quickly in tests.
	origDelay := edgeTransportNodeStatePollDelay
	origInterval := edgeTransportNodeStatePollInterval
	origTimeout := edgeTransportNodeStatePollTimeout
	edgeTransportNodeStatePollDelay = 0
	edgeTransportNodeStatePollInterval = 1
	edgeTransportNodeStatePollTimeout = 10
	defer func() {
		edgeTransportNodeStatePollDelay = origDelay
		edgeTransportNodeStatePollInterval = origInterval
		edgeTransportNodeStatePollTimeout = origTimeout
	}()

	t.Run("Delete success (node disappears after deletion)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTN, restoreTN := setupTransportNodeMock(t, ctrl)
		mockState, restoreState := setupStateClientMock(t, ctrl)
		defer restoreTN()
		defer restoreState()

		mockTN.EXPECT().Delete(tnID, (*bool)(nil), (*bool)(nil)).Return(nil)
		mockState.EXPECT().Get(tnID).Return(mpmodel.TransportNodeState{}, vapiErrors.NotFound{})

		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(tnID)

		err := resourceNsxtEdgeTransportNodeDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtEdgeTransportNodeDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Delete fails when Delete API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTN, restoreTN := setupTransportNodeMock(t, ctrl)
		defer restoreTN()

		mockTN.EXPECT().Delete(tnID, (*bool)(nil), (*bool)(nil)).Return(errors.New("delete API error"))

		res := resourceNsxtEdgeTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(tnID)

		err := resourceNsxtEdgeTransportNodeDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete API error")
	})
}
