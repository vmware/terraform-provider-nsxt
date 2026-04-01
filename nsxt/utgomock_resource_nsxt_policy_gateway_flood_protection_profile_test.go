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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	gwFppID          = "gw-fpp-1"
	gwFppDisplayName = "Test Gateway Flood Protection Profile"
	gwFppDescription = "Test gateway FPP"
	gwFppRevision    = int64(1)
	gwFppPath        = "/infra/flood-protection-profiles/gw-fpp-1"
)

func gwFppStructValue() *data.StructValue {
	profilePath := gwFppPath
	obj := nsxModel.GatewayFloodProtectionProfile{
		DisplayName:  &gwFppDisplayName,
		Description:  &gwFppDescription,
		Revision:     &gwFppRevision,
		Path:         &profilePath,
		ResourceType: nsxModel.FloodProtectionProfile_RESOURCE_TYPE_GATEWAYFLOODPROTECTIONPROFILE,
	}
	converter := bindings.NewTypeConverter()
	profileValue, _ := converter.ConvertToVapi(obj, nsxModel.GatewayFloodProtectionProfileBindingType())
	return profileValue.(*data.StructValue)
}

func minimalGwFppData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":             gwFppDisplayName,
		"description":              gwFppDescription,
		"nsx_id":                   gwFppID,
		"icmp_active_flow_limit":   0,
		"other_active_conn_limit":  0,
		"tcp_half_open_conn_limit": 0,
		"udp_active_flow_limit":    0,
		"nat_active_conn_limit":    0,
	}
}

func setupGwFppMock(t *testing.T, ctrl *gomock.Controller) (*inframocks.MockFloodProtectionProfilesClient, func()) {
	mockSDK := inframocks.NewMockFloodProtectionProfilesClient(ctrl)
	mockWrapper := &cliinfra.StructValueClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliFloodProtectionProfilesClient
	cliFloodProtectionProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.StructValueClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliFloodProtectionProfilesClient = original }
}

func TestMockResourceNsxtPolicyGatewayFloodProtectionProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGwFppMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gwFppID).Return(nil, notFoundErr),
			mockSDK.EXPECT().Patch(gwFppID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gwFppID).Return(gwFppStructValue(), nil),
		)

		res := resourceNsxtPolicyGatewayFloodProtectionProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwFppData())

		err := resourceNsxtPolicyGatewayFloodProtectionProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, gwFppID, d.Id())
	})
}

func TestMockResourceNsxtPolicyGatewayFloodProtectionProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGwFppMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(gwFppID).Return(gwFppStructValue(), nil)

		res := resourceNsxtPolicyGatewayFloodProtectionProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwFppData())
		d.SetId(gwFppID)

		err := resourceNsxtPolicyGatewayFloodProtectionProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, gwFppDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(gwFppID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGatewayFloodProtectionProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwFppData())
		d.SetId(gwFppID)

		err := resourceNsxtPolicyGatewayFloodProtectionProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayFloodProtectionProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwFppData())

		err := resourceNsxtPolicyGatewayFloodProtectionProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayFloodProtectionProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGwFppMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(gwFppID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gwFppID).Return(gwFppStructValue(), nil),
		)

		res := resourceNsxtPolicyGatewayFloodProtectionProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwFppData())
		d.SetId(gwFppID)

		err := resourceNsxtPolicyGatewayFloodProtectionProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayFloodProtectionProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwFppData())

		err := resourceNsxtPolicyGatewayFloodProtectionProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayFloodProtectionProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGwFppMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(gwFppID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyGatewayFloodProtectionProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwFppData())
		d.SetId(gwFppID)

		err := resourceNsxtPolicyFloodProtectionProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
