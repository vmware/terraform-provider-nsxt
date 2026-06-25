//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/infra/route_controllers/InterfacesClient.go -package=mocks -source=<sdk>/services/nsxt/infra/route_controllers/InterfacesClient.go InterfacesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	rcifacemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers"
)

var (
	rcInterfaceID          = "intf-1"
	rcInterfaceDisplayName = "intf-fooname"
	rcInterfaceDescription = "route controller interface mock"
	rcInterfacePath        = "/infra/route-controllers/rc-1/interfaces/intf-1"
	rcInterfaceRevision    = int64(1)
	rcInterfaceParentPath  = "/infra/route-controllers/rc-1"
	rcInterfaceUrpfMode    = model.RouteControllerInterface_URPF_MODE_STRICT
)

func rcInterfaceAPIResponse() model.RouteControllerInterface {
	ipAddr := "10.0.0.1"
	prefix := int64(24)
	return model.RouteControllerInterface{
		Id:          &rcInterfaceID,
		DisplayName: &rcInterfaceDisplayName,
		Description: &rcInterfaceDescription,
		Path:        &rcInterfacePath,
		Revision:    &rcInterfaceRevision,
		UrpfMode:    &rcInterfaceUrpfMode,
		InterfaceAddress: []model.RouteControllerInterfaceAddress{
			{
				InterfaceSubnet: []model.InterfaceSubnet{
					{IpAddresses: []string{ipAddr}, PrefixLen: &prefix},
				},
			},
		},
	}
}

func minimalRCInterfaceData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":        rcInterfaceDisplayName,
		"description":         rcInterfaceDescription,
		"parent_path":         rcInterfaceParentPath,
		"floating_ip_subnets": []interface{}{"10.0.1.1/24"},
		"interface_address": []interface{}{
			map[string]interface{}{
				"subnets":                        []interface{}{"10.0.0.1/24"},
				"portgroup_id":                   "mock-compute-manager-uuid:dvportgroup-1",
				"virtual_network_appliance_path": "",
			},
		},
	}
}

func setupRCInterfaceMocks(t *testing.T, ctrl *gomock.Controller) (
	*rcifacemocks.MockInterfacesClient,
	func(),
) {
	mockIfaceSDK := rcifacemocks.NewMockInterfacesClient(ctrl)

	ifaceWrapper := &cliinfra.RouteControllerInterfaceClientContext{
		Client:     mockIfaceSDK,
		ClientType: utl.Local,
	}

	orig := cliRCInterfaceClient

	cliRCInterfaceClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.RouteControllerInterfaceClientContext {
		return ifaceWrapper
	}

	return mockIfaceSDK, func() {
		cliRCInterfaceClient = orig
	}
}

func TestMockResourceNsxtPolicyRouteControllerInterfaceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIfaceSDK, restore := setupRCInterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		resp := rcInterfaceAPIResponse()
		gomock.InOrder(
			mockIfaceSDK.EXPECT().Patch(routeControllerID, gomock.Any(), gomock.Any()).Return(nil),
			mockIfaceSDK.EXPECT().Get(routeControllerID, gomock.Any()).Return(resp, nil),
		)

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRCInterfaceData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockIfaceSDK.EXPECT().Get(routeControllerID, rcInterfaceID).Return(rcInterfaceAPIResponse(), nil)

		res := resourceNsxtPolicyRouteControllerInterface()
		data := minimalRCInterfaceData()
		data["nsx_id"] = rcInterfaceID
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when API returns error", func(t *testing.T) {
		mockIfaceSDK.EXPECT().Patch(routeControllerID, gomock.Any(), gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRCInterfaceData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtPolicyRouteControllerInterfaceRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIfaceSDK, restore := setupRCInterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockIfaceSDK.EXPECT().Get(routeControllerID, rcInterfaceID).Return(rcInterfaceAPIResponse(), nil)

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcInterfaceParentPath,
		})
		d.SetId(rcInterfaceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcInterfaceDisplayName, d.Get("display_name"))
		assert.Equal(t, rcInterfaceDescription, d.Get("description"))
		assert.Equal(t, rcInterfacePath, d.Get("path"))
		assert.Equal(t, int(rcInterfaceRevision), d.Get("revision"))
		assert.Equal(t, rcInterfaceUrpfMode, d.Get("urpf_mode"))
		assert.Equal(t, rcInterfaceID, d.Id())
	})

	t.Run("Read fails when not found", func(t *testing.T) {
		mockIfaceSDK.EXPECT().Get(routeControllerID, rcInterfaceID).Return(model.RouteControllerInterface{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcInterfaceParentPath,
		})
		d.SetId(rcInterfaceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcInterfaceParentPath,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RouteControllerInterface ID")
	})
}

func TestMockResourceNsxtPolicyRouteControllerInterfaceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIfaceSDK, restore := setupRCInterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		resp := rcInterfaceAPIResponse()
		gomock.InOrder(
			mockIfaceSDK.EXPECT().Patch(routeControllerID, rcInterfaceID, gomock.Any()).Return(nil),
			mockIfaceSDK.EXPECT().Get(routeControllerID, rcInterfaceID).Return(resp, nil),
		)

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRCInterfaceData())
		d.SetId(rcInterfaceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRCInterfaceData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RouteControllerInterface ID")
	})

	t.Run("Update fails when API returns error", func(t *testing.T) {
		mockIfaceSDK.EXPECT().Patch(routeControllerID, rcInterfaceID, gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRCInterfaceData())
		d.SetId(rcInterfaceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtPolicyRouteControllerInterfaceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIfaceSDK, restore := setupRCInterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockIfaceSDK.EXPECT().Delete(routeControllerID, rcInterfaceID).Return(nil)

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcInterfaceParentPath,
		})
		d.SetId(rcInterfaceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcInterfaceParentPath,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RouteControllerInterface ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockIfaceSDK.EXPECT().Delete(routeControllerID, rcInterfaceID).Return(errors.New("API error"))

		res := resourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcInterfaceParentPath,
		})
		d.SetId(rcInterfaceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerInterfaceDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}
