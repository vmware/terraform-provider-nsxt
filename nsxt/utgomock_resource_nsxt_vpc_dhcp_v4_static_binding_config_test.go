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

	apisubnets "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs/subnets"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	subnetsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs/subnets"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dhcpBindingID          = "dhcp-binding-id"
	dhcpBindingDisplayName = "test-dhcp-binding"
	dhcpBindingDescription = "Test DHCP v4 Static Binding"
	dhcpBindingRevision    = int64(1)
	dhcpBindingParentPath  = "/orgs/default/projects/project1/vpcs/vpc1/subnets/subnet1"
	dhcpBindingMacAddr     = "AA:BB:CC:DD:EE:FF"
	dhcpBindingIPAddr      = "10.0.0.10"
)

func dhcpV4BindingStructValue() *data.StructValue {
	displayName := dhcpBindingDisplayName
	description := dhcpBindingDescription
	revision := dhcpBindingRevision
	macAddr := dhcpBindingMacAddr
	ipAddr := dhcpBindingIPAddr
	obj := nsxModel.DhcpV4StaticBindingConfig{
		Id:           &dhcpBindingID,
		DisplayName:  &displayName,
		Description:  &description,
		Revision:     &revision,
		ResourceType: nsxModel.DhcpStaticBindingConfig_RESOURCE_TYPE_DHCPV4STATICBINDINGCONFIG,
		MacAddress:   &macAddr,
		IpAddress:    &ipAddr,
	}
	converter := bindings.NewTypeConverter()
	structVal, errs := converter.ConvertToVapi(obj, nsxModel.DhcpV4StaticBindingConfigBindingType())
	if errs != nil {
		panic(errs[0])
	}
	return structVal.(*data.StructValue)
}

func minimalDhcpBindingData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": dhcpBindingDisplayName,
		"description":  dhcpBindingDescription,
		"nsx_id":       dhcpBindingID,
		"parent_path":  dhcpBindingParentPath,
		"mac_address":  dhcpBindingMacAddr,
		"ip_address":   dhcpBindingIPAddr,
	}
}

func setupDhcpBindingMock(t *testing.T, ctrl *gomock.Controller) (*subnetsmocks.MockDhcpStaticBindingConfigsClient, func()) {
	mockSDK := subnetsmocks.NewMockDhcpStaticBindingConfigsClient(ctrl)
	mockWrapper := &apisubnets.VpcDhcpStaticBindingClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	orig := cliVpcSubnetDhcpStaticBindingConfigsClient
	cliVpcSubnetDhcpStaticBindingConfigsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apisubnets.VpcDhcpStaticBindingClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcSubnetDhcpStaticBindingConfigsClient = orig }
}

func TestMockResourceNsxtVpcDhcpV4StaticBindingCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpBindingMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpBindingID).Return(nil, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpBindingID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpBindingID).Return(dhcpV4BindingStructValue(), nil),
		)

		res := resourceNsxtVpcSubnetDhcpV4StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpBindingData())

		err := resourceNsxtVpcSubnetDhcpV4StaticBindingConfigCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dhcpBindingID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcSubnetDhcpV4StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpBindingData())

		err := resourceNsxtVpcSubnetDhcpV4StaticBindingConfigCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcDhcpV4StaticBindingRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpBindingMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpBindingID).Return(dhcpV4BindingStructValue(), nil)

		res := resourceNsxtVpcSubnetDhcpV4StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpBindingData())
		d.SetId(dhcpBindingID)

		err := resourceNsxtVpcSubnetDhcpV4StaticBindingConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dhcpBindingDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpBindingMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpBindingID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtVpcSubnetDhcpV4StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpBindingData())
		d.SetId(dhcpBindingID)

		err := resourceNsxtVpcSubnetDhcpV4StaticBindingConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcDhcpV4StaticBindingUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpBindingMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpBindingID, gomock.Any()).Return(dhcpV4BindingStructValue(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpBindingID).Return(dhcpV4BindingStructValue(), nil),
		)

		res := resourceNsxtVpcSubnetDhcpV4StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpBindingData())
		d.SetId(dhcpBindingID)

		err := resourceNsxtVpcSubnetDhcpV4StaticBindingConfigUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcDhcpV4StaticBindingDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpBindingMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpBindingID).Return(nil)

		res := resourceNsxtVpcSubnetDhcpV4StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpBindingData())
		d.SetId(dhcpBindingID)

		err := resourceNsxtVpcSubnetDhcpV4StaticBindingConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
