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

	tier0sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	t0lsapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
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
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

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
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

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
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

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
