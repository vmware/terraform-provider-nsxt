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
	dhcpV6BindingID          = "dhcp-v6-binding-id"
	dhcpV6BindingDisplayName = "test-dhcp-v6-binding"
	dhcpV6BindingDescription = "Test DHCP v6 Static Binding"
	dhcpV6BindingRevision    = int64(1)
	dhcpV6BindingParentPath  = "/orgs/default/projects/project1/vpcs/vpc1/subnets/subnet1"
	dhcpV6BindingMacAddr     = "AA:BB:CC:DD:EE:FF"
	dhcpV6BindingIPAddr      = "2001:db8::10"
)

func dhcpV6BindingStructValue() *data.StructValue {
	displayName := dhcpV6BindingDisplayName
	description := dhcpV6BindingDescription
	revision := dhcpV6BindingRevision
	macAddr := dhcpV6BindingMacAddr
	lease := int64(86400)
	obj := nsxModel.DhcpV6StaticBindingConfig{
		Id:           &dhcpV6BindingID,
		DisplayName:  &displayName,
		Description:  &description,
		Revision:     &revision,
		ResourceType: nsxModel.DhcpStaticBindingConfig_RESOURCE_TYPE_DHCPV6STATICBINDINGCONFIG,
		MacAddress:   &macAddr,
		LeaseTime:    &lease,
		IpAddresses:  []string{dhcpV6BindingIPAddr},
	}
	converter := bindings.NewTypeConverter()
	structVal, errs := converter.ConvertToVapi(obj, nsxModel.DhcpV6StaticBindingConfigBindingType())
	if errs != nil {
		panic(errs[0])
	}
	return structVal.(*data.StructValue)
}

func minimalDhcpV6BindingData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": dhcpV6BindingDisplayName,
		"description":  dhcpV6BindingDescription,
		"nsx_id":       dhcpV6BindingID,
		"parent_path":  dhcpV6BindingParentPath,
		"mac_address":  dhcpV6BindingMacAddr,
		"ip_addresses": []interface{}{dhcpV6BindingIPAddr},
	}
}

func setupDhcpV6BindingMock(t *testing.T, ctrl *gomock.Controller) (*subnetsmocks.MockDhcpStaticBindingConfigsClient, func()) {
	mockSDK := subnetsmocks.NewMockDhcpStaticBindingConfigsClient(ctrl)
	mockWrapper := &apisubnets.VpcDhcpStaticBindingClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	orig := cliVpcSubnetDhcpV6StaticBindingConfigsClient
	cliVpcSubnetDhcpV6StaticBindingConfigsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apisubnets.VpcDhcpStaticBindingClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcSubnetDhcpV6StaticBindingConfigsClient = orig }
}

func TestMockResourceNsxtVpcDhcpV6StaticBindingCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpV6BindingMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpV6BindingID).Return(nil, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpV6BindingID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpV6BindingID).Return(dhcpV6BindingStructValue(), nil),
		)

		res := resourceNsxtVpcSubnetDhcpV6StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpV6BindingData())

		err := resourceNsxtVpcSubnetDhcpV6StaticBindingConfigCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dhcpV6BindingID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcSubnetDhcpV6StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpV6BindingData())

		err := resourceNsxtVpcSubnetDhcpV6StaticBindingConfigCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcDhcpV6StaticBindingRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpV6BindingMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpV6BindingID).Return(dhcpV6BindingStructValue(), nil)

		res := resourceNsxtVpcSubnetDhcpV6StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpV6BindingData())
		d.SetId(dhcpV6BindingID)

		err := resourceNsxtVpcSubnetDhcpV6StaticBindingConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dhcpV6BindingDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpV6BindingMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpV6BindingID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtVpcSubnetDhcpV6StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpV6BindingData())
		d.SetId(dhcpV6BindingID)

		err := resourceNsxtVpcSubnetDhcpV6StaticBindingConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcDhcpV6StaticBindingUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpV6BindingMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpV6BindingID, gomock.Any()).Return(dhcpV6BindingStructValue(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpV6BindingID).Return(dhcpV6BindingStructValue(), nil),
		)

		res := resourceNsxtVpcSubnetDhcpV6StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpV6BindingData())
		d.SetId(dhcpV6BindingID)

		err := resourceNsxtVpcSubnetDhcpV6StaticBindingConfigUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcDhcpV6StaticBindingDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpV6BindingMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dhcpV6BindingID).Return(nil)

		res := resourceNsxtVpcSubnetDhcpV6StaticBindingConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpV6BindingData())
		d.SetId(dhcpV6BindingID)

		err := resourceNsxtVpcSubnetDhcpV6StaticBindingConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
