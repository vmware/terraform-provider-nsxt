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

	tier1sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s"
	t1lsapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/locale_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t1lsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
	t1intmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/locale_services"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	t1IntID              = "t1-intf-1"
	t1IntDisplayName     = "Test T1 Interface"
	t1IntDescription     = "Test Tier1 interface"
	t1IntRevision        = int64(1)
	t1IntGwPath          = "/infra/tier-1s/t1-gw-1"
	t1IntGwID            = "t1-gw-1"
	t1IntLocaleServiceID = "default"
	t1IntSubnet          = "192.168.20.1/24"
	t1IntSegmentPath     = "/infra/segments/seg-2"
)

func t1InterfaceAPIResponse() nsxModel.Tier1Interface {
	return nsxModel.Tier1Interface{
		Id:          &t1IntID,
		DisplayName: &t1IntDisplayName,
		Description: &t1IntDescription,
		Revision:    &t1IntRevision,
	}
}

func minimalT1InterfaceData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": t1IntDisplayName,
		"description":  t1IntDescription,
		"gateway_path": t1IntGwPath,
		"segment_path": t1IntSegmentPath,
		"subnets":      []interface{}{t1IntSubnet},
	}
}

func t1InterfaceDataWithLocaleService() map[string]interface{} {
	d := minimalT1InterfaceData()
	d["locale_service_id"] = t1IntLocaleServiceID
	return d
}

func setupT1InterfaceMocks(t *testing.T, ctrl *gomock.Controller) (
	*t1intmocks.MockInterfacesClient,
	*t1lsmocks.MockLocaleServicesClient,
	func(),
) {
	mockIntfSDK := t1intmocks.NewMockInterfacesClient(ctrl)
	mockLSSDK := t1lsmocks.NewMockLocaleServicesClient(ctrl)

	intfWrapper := &t1lsapi.Tier1InterfaceClientContext{
		Client:     mockIntfSDK,
		ClientType: utl.Local,
	}
	lsWrapper := &tier1sapi.LocaleServicesClientContext{
		Client:     mockLSSDK,
		ClientType: utl.Local,
	}

	origIntf := cliTier1InterfacesClient
	origLS := cliTier1LocaleServicesClient
	cliTier1InterfacesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t1lsapi.Tier1InterfaceClientContext {
		return intfWrapper
	}
	cliTier1LocaleServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier1sapi.LocaleServicesClientContext {
		return lsWrapper
	}

	return mockIntfSDK, mockLSSDK, func() {
		cliTier1InterfacesClient = origIntf
		cliTier1LocaleServicesClient = origLS
	}
}

func TestMockResourceNsxtPolicyTier1GatewayInterfaceCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIntfSDK, mockLSSDK, restore := setupT1InterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		lsID := t1IntLocaleServiceID
		lsResponse := nsxModel.LocaleServices{Id: &lsID}
		gomock.InOrder(
			mockLSSDK.EXPECT().Get(t1IntGwID, defaultPolicyLocaleServiceID).Return(lsResponse, nil),
			mockIntfSDK.EXPECT().Patch(t1IntGwID, t1IntLocaleServiceID, gomock.Any(), gomock.Any()).Return(nil),
			mockIntfSDK.EXPECT().Get(t1IntGwID, t1IntLocaleServiceID, gomock.Any()).Return(t1InterfaceAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTier1GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1InterfaceData())

		err := resourceNsxtPolicyTier1GatewayInterfaceCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, t1IntDisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyTier1GatewayInterfaceRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIntfSDK, _, restore := setupT1InterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockIntfSDK.EXPECT().Get(t1IntGwID, t1IntLocaleServiceID, t1IntID).Return(t1InterfaceAPIResponse(), nil)

		res := resourceNsxtPolicyTier1GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t1InterfaceDataWithLocaleService())
		d.SetId(t1IntID)

		err := resourceNsxtPolicyTier1GatewayInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, t1IntDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockIntfSDK.EXPECT().Get(t1IntGwID, t1IntLocaleServiceID, t1IntID).Return(nsxModel.Tier1Interface{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTier1GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t1InterfaceDataWithLocaleService())
		d.SetId(t1IntID)

		err := resourceNsxtPolicyTier1GatewayInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier1GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t1InterfaceDataWithLocaleService())

		err := resourceNsxtPolicyTier1GatewayInterfaceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier1GatewayInterfaceUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIntfSDK, _, restore := setupT1InterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockIntfSDK.EXPECT().Update(t1IntGwID, t1IntLocaleServiceID, t1IntID, gomock.Any()).Return(t1InterfaceAPIResponse(), nil),
			mockIntfSDK.EXPECT().Get(t1IntGwID, t1IntLocaleServiceID, t1IntID).Return(t1InterfaceAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTier1GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t1InterfaceDataWithLocaleService())
		d.SetId(t1IntID)

		err := resourceNsxtPolicyTier1GatewayInterfaceUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when IDs are empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier1GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1InterfaceData())

		err := resourceNsxtPolicyTier1GatewayInterfaceUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier1GatewayInterfaceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIntfSDK, _, restore := setupT1InterfaceMocks(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockIntfSDK.EXPECT().Delete(t1IntGwID, t1IntLocaleServiceID, t1IntID).Return(nil)

		res := resourceNsxtPolicyTier1GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, t1InterfaceDataWithLocaleService())
		d.SetId(t1IntID)

		err := resourceNsxtPolicyTier1GatewayInterfaceDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when IDs are empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier1GatewayInterface()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1InterfaceData())

		err := resourceNsxtPolicyTier1GatewayInterfaceDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
