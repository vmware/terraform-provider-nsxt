// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apisubnets "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs/subnets"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	subnetsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs/subnets"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	extAddrPortPath  = "/orgs/default/projects/project1/vpcs/vpc1/subnets/subnet1/ports/port1"
	extAddrAllocPath = "/orgs/default/projects/project1/vpcs/vpc1/subnets/subnet1/ip-allocations/alloc1"
	extAddrExtIP     = "192.168.1.100"
)

func extAddrPortAPIResponse() nsxModel.VpcSubnetPort {
	allocPath := extAddrAllocPath
	extIP := extAddrExtIP
	return nsxModel.VpcSubnetPort{
		Id: func() *string { s := "port1"; return &s }(),
		ExternalAddressBinding: &nsxModel.ExternalAddressBinding{
			AllocatedExternalIpPath: &allocPath,
			ExternalIpAddress:       &extIP,
		},
	}
}

func portWithNoExternalAddress() nsxModel.VpcSubnetPort {
	return nsxModel.VpcSubnetPort{
		Id: func() *string { s := "port1"; return &s }(),
	}
}

func setupExtAddrMock(t *testing.T, ctrl *gomock.Controller) (*subnetsmocks.MockPortsClient, func()) {
	mockPortsSDK := subnetsmocks.NewMockPortsClient(ctrl)
	mockWrapper := &apisubnets.VpcSubnetPortClientContext{
		Client:     mockPortsSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	orig := cliVpcSubnetPortsClient
	cliVpcSubnetPortsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apisubnets.VpcSubnetPortClientContext {
		return mockWrapper
	}
	return mockPortsSDK, func() { cliVpcSubnetPortsClient = orig }
}

func TestMockResourceNsxtVpcExternalAddressCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPortsSDK, restore := setupExtAddrMock(t, ctrl)
		defer restore()

		portWithNoAddr := portWithNoExternalAddress()
		gomock.InOrder(
			// updatePort: Get then Update
			mockPortsSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(portWithNoAddr, nil),
			mockPortsSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(extAddrPortAPIResponse(), nil),
			// Read after create
			mockPortsSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(extAddrPortAPIResponse(), nil),
		)

		res := resourceNsxtVpcExternalAddress()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path":                extAddrPortPath,
			"allocated_external_ip_path": extAddrAllocPath,
		})

		err := resourceNsxtVpcExternalAddressCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, extAddrAllocPath, d.Get("allocated_external_ip_path"))
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcExternalAddress()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path":                extAddrPortPath,
			"allocated_external_ip_path": extAddrAllocPath,
		})

		err := resourceNsxtVpcExternalAddressCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create fails when port Get returns error", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPortsSDK, restore := setupExtAddrMock(t, ctrl)
		defer restore()

		mockPortsSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.VpcSubnetPort{}, errors.New("get failed"))

		res := resourceNsxtVpcExternalAddress()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path":                extAddrPortPath,
			"allocated_external_ip_path": extAddrAllocPath,
		})

		err := resourceNsxtVpcExternalAddressCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcExternalAddressRead(t *testing.T) {
	t.Run("Read success with external address", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPortsSDK, restore := setupExtAddrMock(t, ctrl)
		defer restore()

		mockPortsSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(extAddrPortAPIResponse(), nil)

		res := resourceNsxtVpcExternalAddress()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path":                extAddrPortPath,
			"allocated_external_ip_path": extAddrAllocPath,
		})
		d.SetId("some-uuid")

		err := resourceNsxtVpcExternalAddressRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, extAddrAllocPath, d.Get("allocated_external_ip_path"))
		assert.Equal(t, extAddrExtIP, d.Get("external_ip_address"))
	})

	t.Run("Read with no external address binding", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPortsSDK, restore := setupExtAddrMock(t, ctrl)
		defer restore()

		mockPortsSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(portWithNoExternalAddress(), nil)

		res := resourceNsxtVpcExternalAddress()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path":                extAddrPortPath,
			"allocated_external_ip_path": "",
		})
		d.SetId("some-uuid")

		err := resourceNsxtVpcExternalAddressRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Get("allocated_external_ip_path"))
	})
}

func TestMockResourceNsxtVpcExternalAddressDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPortsSDK, restore := setupExtAddrMock(t, ctrl)
		defer restore()

		portWithAddr := extAddrPortAPIResponse()
		gomock.InOrder(
			mockPortsSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(portWithAddr, nil),
			mockPortsSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(portWithNoExternalAddress(), nil),
		)

		res := resourceNsxtVpcExternalAddress()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path":                extAddrPortPath,
			"allocated_external_ip_path": extAddrAllocPath,
		})
		d.SetId("some-uuid")

		err := resourceNsxtVpcExternalAddressDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
