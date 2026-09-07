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

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	tier0sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	t0lsapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0sdkmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	t0lsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	t0intmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	t0IntID              = "t0-intf-1"
	t0IntDisplayName     = "Test T0 Interface"
	t0IntDescription     = "Test Tier0 interface"
	t0IntRevision        = int64(1)
	t0IntGwPath          = "/infra/tier-0s/t0-gw-1"
	t0IntGwID            = "t0-gw-1"
	t0IntLocaleServiceID = "default"
	t0IntSubnet          = "192.168.10.1/24"
	t0IntSegmentPath     = "/infra/segments/seg-1"
)

func t0InterfaceAPIResponse() nsxModel.Tier0Interface {
	ifType := nsxModel.Tier0Interface_TYPE_EXTERNAL
	return nsxModel.Tier0Interface{
		Id:          &t0IntID,
		DisplayName: &t0IntDisplayName,
		Description: &t0IntDescription,
		Revision:    &t0IntRevision,
		Type_:       &ifType,
	}
}

func minimalT0InterfaceData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": t0IntDisplayName,
		"description":  t0IntDescription,
		"gateway_path": t0IntGwPath,
		"type":         nsxModel.Tier0Interface_TYPE_EXTERNAL,
		"segment_path": t0IntSegmentPath,
		"subnets":      []interface{}{t0IntSubnet},
		"urpf_mode":    nsxModel.Tier0Interface_URPF_MODE_STRICT,
		"enable_pim":   false,
	}
}

func t0InterfaceDataWithLocaleService() map[string]interface{} {
	d := minimalT0InterfaceData()
	d["locale_service_id"] = t0IntLocaleServiceID
	return d
}

func setupT0InterfaceMocks(t *testing.T, ctrl *gomock.Controller) (
	*t0intmocks.MockInterfacesClient,
	*t0lsmocks.MockLocaleServicesClient,
	func(),
) {
	mockIntfSDK := t0intmocks.NewMockInterfacesClient(ctrl)
	mockIntfWrapper := &t0lsapi.Tier0InterfaceClientContext{
		Client:     mockIntfSDK,
		ClientType: utl.Local,
	}
	mockLSSDK := t0lsmocks.NewMockLocaleServicesClient(ctrl)
	mockLSWrapper := &tier0sapi.LocaleServicesClientContext{
		Client:     mockLSSDK,
		ClientType: utl.Local,
	}

	origIntf := cliTier0InterfacesClient
	origLS := cliTier0LocaleServicesClient
	cliTier0InterfacesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t0lsapi.Tier0InterfaceClientContext {
		return mockIntfWrapper
	}
	cliTier0LocaleServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0sapi.LocaleServicesClientContext {
		return mockLSWrapper
	}

	return mockIntfSDK, mockLSSDK, func() {
		cliTier0InterfacesClient = origIntf
		cliTier0LocaleServicesClient = origLS
	}
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceCreate(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIntfSDK, mockLSSDK, restore := setupT0InterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		lsID := t0IntLocaleServiceID
		lsResponse := nsxModel.LocaleServices{Id: &lsID}
		gomock.InOrder(
			mockLSSDK.EXPECT().Get(t0IntGwID, defaultPolicyLocaleServiceID).Return(lsResponse, nil),
			mockIntfSDK.EXPECT().Patch(t0IntGwID, t0IntLocaleServiceID, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
			mockIntfSDK.EXPECT().Get(t0IntGwID, t0IntLocaleServiceID, gomock.Any()).Return(t0InterfaceAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT0InterfaceData())

		err := resourceNsxtPolicyTier0GatewayInterfaceCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, t0IntDisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceRead(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIntfSDK, _, restore := setupT0InterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockIntfSDK.EXPECT().Get(t0IntGwID, t0IntLocaleServiceID, t0IntID).Return(t0InterfaceAPIResponse(), nil)

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t0InterfaceDataWithLocaleService())
		d.SetId(t0IntID)

		err := resourceNsxtPolicyTier0GatewayInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, t0IntDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockIntfSDK.EXPECT().Get(t0IntGwID, t0IntLocaleServiceID, t0IntID).Return(nsxModel.Tier0Interface{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t0InterfaceDataWithLocaleService())
		d.SetId(t0IntID)

		err := resourceNsxtPolicyTier0GatewayInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t0InterfaceDataWithLocaleService())

		err := resourceNsxtPolicyTier0GatewayInterfaceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceUpdate(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIntfSDK, _, restore := setupT0InterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockIntfSDK.EXPECT().Update(t0IntGwID, t0IntLocaleServiceID, t0IntID, gomock.Any(), gomock.Any()).Return(t0InterfaceAPIResponse(), nil),
			mockIntfSDK.EXPECT().Get(t0IntGwID, t0IntLocaleServiceID, t0IntID).Return(t0InterfaceAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t0InterfaceDataWithLocaleService())
		d.SetId(t0IntID)

		err := resourceNsxtPolicyTier0GatewayInterfaceUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when IDs are empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT0InterfaceData())

		err := resourceNsxtPolicyTier0GatewayInterfaceUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIntfSDK, _, restore := setupT0InterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockIntfSDK.EXPECT().Delete(t0IntGwID, t0IntLocaleServiceID, t0IntID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t0InterfaceDataWithLocaleService())
		d.SetId(t0IntID)

		err := resourceNsxtPolicyTier0GatewayInterfaceDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when IDs are empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT0InterfaceData())

		err := resourceNsxtPolicyTier0GatewayInterfaceDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayInterface_subnetPathVersion(t *testing.T) {
	t.Run("Create fails with subnet_path before NSX 9.2.0", func(t *testing.T) {

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0IntDisplayName,
			"description":  t0IntDescription,
			"gateway_path": t0IntGwPath,
			"type":         nsxModel.Tier0Interface_TYPE_EXTERNAL,
			"subnet_path":  "/orgs/default/projects/p1/vpcs/v1/subnets/s1",
			"subnets":      []interface{}{t0IntSubnet},
			"urpf_mode":    nsxModel.Tier0Interface_URPF_MODE_STRICT,
			"enable_pim":   false,
		})

		err := resourceNsxtPolicyTier0GatewayInterfaceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceRead_subnetPathVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIntfSDK, _, restore := setupT0InterfaceMocks(t, ctrl)
	defer restore()

	sp := "/orgs/default/projects/p1/vpcs/v1/subnets/s1"

	t.Run("subnet_path omitted in state before 9.2.0", func(t *testing.T) {

		iface := t0InterfaceAPIResponse()
		iface.SubnetPath = &sp
		mockIntfSDK.EXPECT().Get(t0IntGwID, t0IntLocaleServiceID, t0IntID).Return(iface, nil)

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t0InterfaceDataWithLocaleService())
		d.SetId(t0IntID)

		err := resourceNsxtPolicyTier0GatewayInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Get("subnet_path"))
	})

	t.Run("subnet_path set in state on NSX 9.2.0", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		iface := t0InterfaceAPIResponse()
		iface.SubnetPath = &sp
		mockIntfSDK.EXPECT().Get(t0IntGwID, t0IntLocaleServiceID, t0IntID).Return(iface, nil)

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t0InterfaceDataWithLocaleService())
		d.SetId(t0IntID)

		err := resourceNsxtPolicyTier0GatewayInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, sp, d.Get("subnet_path"))
	})
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceCreateValidation(t *testing.T) {
	t.Run("fails when neither segment_path nor subnet_path is set for non-loopback type", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0IntDisplayName,
			"gateway_path": t0IntGwPath,
			"type":         nsxModel.Tier0Interface_TYPE_EXTERNAL,
		})

		err := resourceNsxtPolicyTier0GatewayInterfaceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mandatory")
	})

	t.Run("fails when enable_pim is set on Global Manager", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0IntDisplayName,
			"gateway_path": t0IntGwPath,
			"type":         nsxModel.Tier0Interface_TYPE_EXTERNAL,
			"segment_path": t0IntSegmentPath,
			"enable_pim":   true,
		})

		err := resourceNsxtPolicyTier0GatewayInterfaceCreate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "enable_pim")
	})

	t.Run("fails when site_path is missing on Global Manager", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0IntDisplayName,
			"gateway_path": t0IntGwPath,
			"type":         nsxModel.Tier0Interface_TYPE_EXTERNAL,
			"segment_path": t0IntSegmentPath,
			"enable_pim":   false,
		})

		err := resourceNsxtPolicyTier0GatewayInterfaceCreate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "site_path")
	})

	t.Run("fails when site_path is set on Local Manager", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0IntDisplayName,
			"gateway_path": t0IntGwPath,
			"type":         nsxModel.Tier0Interface_TYPE_EXTERNAL,
			"segment_path": t0IntSegmentPath,
			"site_path":    "/infra/sites/default",
		})

		err := resourceNsxtPolicyTier0GatewayInterfaceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceCreateExistsCheck(t *testing.T) {
	t.Run("fails when interface with nsx_id already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockIntfSDK, mockLSSDK, restore := setupT0InterfaceMocks(t, ctrl)
		defer restore()
		mockLSSDK.EXPECT().Get(t0IntGwID, defaultPolicyLocaleServiceID).Return(nsxModel.LocaleServices{Id: &t0IntLocaleServiceID}, nil)
		mockIntfSDK.EXPECT().Get(t0IntGwID, t0IntLocaleServiceID, t0IntID).Return(t0InterfaceAPIResponse(), nil)

		data := minimalT0InterfaceData()
		data["nsx_id"] = t0IntID
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTier0GatewayInterfaceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("propagates a non-not-found error from the exists check", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockIntfSDK, mockLSSDK, restore := setupT0InterfaceMocks(t, ctrl)
		defer restore()
		mockLSSDK.EXPECT().Get(t0IntGwID, defaultPolicyLocaleServiceID).Return(nsxModel.LocaleServices{Id: &t0IntLocaleServiceID}, nil)
		mockIntfSDK.EXPECT().Get(t0IntGwID, t0IntLocaleServiceID, t0IntID).Return(nsxModel.Tier0Interface{}, vapiErrors.InternalServerError{})

		data := minimalT0InterfaceData()
		data["nsx_id"] = t0IntID
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTier0GatewayInterfaceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceCreateEdgeClusterMissing(t *testing.T) {
	t.Run("fails when no locale service with an edge cluster is found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockLSSDK, restore := setupT0InterfaceMocks(t, ctrl)
		defer restore()
		mockLSSDK.EXPECT().Get(t0IntGwID, defaultPolicyLocaleServiceID).Return(nsxModel.LocaleServices{}, vapiErrors.NotFound{})
		mockLSSDK.EXPECT().List(t0IntGwID, (*string)(nil), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.LocaleServicesListResult{}, vapiErrors.NotFound{},
		)

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT0InterfaceData())

		err := resourceNsxtPolicyTier0GatewayInterfaceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Edge cluster is mandatory")
	})
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceCreateOspf(t *testing.T) {
	t.Run("fails when ospf is set on a non-EXTERNAL interface", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockLSSDK, restore := setupT0InterfaceMocks(t, ctrl)
		defer restore()
		mockLSSDK.EXPECT().Get(t0IntGwID, defaultPolicyLocaleServiceID).Return(nsxModel.LocaleServices{Id: &t0IntLocaleServiceID}, nil)

		data := map[string]interface{}{
			"display_name": t0IntDisplayName,
			"gateway_path": t0IntGwPath,
			"type":         nsxModel.Tier0Interface_TYPE_LOOPBACK,
			"urpf_mode":    nsxModel.Tier0Interface_URPF_MODE_STRICT,
			"ospf": []interface{}{
				map[string]interface{}{
					"enabled":          true,
					"area_path":        "/infra/tier-0s/t0-gw-1/locale-services/default/ospf/areas/0.0.0.0",
					"enable_bfd":       false,
					"hello_interval":   10,
					"dead_interval":    40,
					"network_type":     nsxModel.PolicyInterfaceOspfConfig_NETWORK_TYPE_BROADCAST,
					"bfd_profile_path": "",
				},
			},
		}
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTier0GatewayInterfaceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Ospf")
	})
}

func setupT0GatewayImportMock(t *testing.T, ctrl *gomock.Controller) (*t0sdkmocks.MockTier0sClient, func()) {
	mockTier0sSDK := t0sdkmocks.NewMockTier0sClient(ctrl)
	tier0Wrapper := &cliinfra.Tier0ClientContext{
		Client:     mockTier0sSDK,
		ClientType: utl.Local,
	}
	original := cliTier0sClient
	cliTier0sClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}
	return mockTier0sSDK, func() { cliTier0sClient = original }
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceImport(t *testing.T) {
	t.Run("succeeds with a valid composite ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTier0sSDK, restore := setupT0GatewayImportMock(t, ctrl)
		defer restore()
		gwPath := t0IntGwPath
		mockTier0sSDK.EXPECT().Get(t0IntGwID).Return(nsxModel.Tier0{Path: &gwPath}, nil)

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(t0IntGwID + "/" + t0IntLocaleServiceID + "/" + t0IntID)

		out, err := resourceNsxtPolicyTier0GatewayInterfaceImport(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, t0IntID, out[0].Id())
		assert.Equal(t, t0IntGwPath, out[0].Get("gateway_path"))
		assert.Equal(t, t0IntLocaleServiceID, out[0].Get("locale_service_id"))
	})

	t.Run("fails with a malformed ID", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("not-enough-parts")

		_, err := resourceNsxtPolicyTier0GatewayInterfaceImport(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("propagates a gateway lookup error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockTier0sSDK, restore := setupT0GatewayImportMock(t, ctrl)
		defer restore()
		mockTier0sSDK.EXPECT().Get(t0IntGwID).Return(nsxModel.Tier0{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(t0IntGwID + "/" + t0IntLocaleServiceID + "/" + t0IntID)

		_, err := resourceNsxtPolicyTier0GatewayInterfaceImport(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayInterfaceReadWithOspf(t *testing.T) {
	t.Run("Read populates the ospf block from the API object", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockIntfSDK, _, restore := setupT0InterfaceMocks(t, ctrl)
		defer restore()

		iface := t0InterfaceAPIResponse()
		enabled, enableBfd := true, false
		areaPath := "/infra/tier-0s/t0-gw-1/locale-services/default/ospf/areas/0.0.0.0"
		networkType := nsxModel.PolicyInterfaceOspfConfig_NETWORK_TYPE_BROADCAST
		hello, dead := int64(10), int64(40)
		iface.Ospf = &nsxModel.PolicyInterfaceOspfConfig{
			Enabled:       &enabled,
			EnableBfd:     &enableBfd,
			OspfArea:      &areaPath,
			NetworkType:   &networkType,
			HelloInterval: &hello,
			DeadInterval:  &dead,
		}
		mockIntfSDK.EXPECT().Get(t0IntGwID, t0IntLocaleServiceID, t0IntID).Return(iface, nil)

		res := resourceNsxtPolicyTier0GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t0InterfaceDataWithLocaleService())
		d.SetId(t0IntID)

		err := resourceNsxtPolicyTier0GatewayInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		ospf := d.Get("ospf").([]interface{})
		require.Len(t, ospf, 1)
		assert.Equal(t, areaPath, ospf[0].(map[string]interface{})["area_path"])
	})
}
