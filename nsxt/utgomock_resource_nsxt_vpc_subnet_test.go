//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apivpcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	apisubnets "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs/subnets"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
	subnetsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs/subnets"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	subnetID          = "subnet-id"
	subnetDisplayName = "test-subnet"
	subnetDescription = "Test VPC Subnet"
	subnetRevision    = int64(1)
	subnetIpBlock     = "/orgs/default/projects/project1/infra/ip-blocks/block1"
)

func subnetAPIResponse() nsxModel.VpcSubnet {
	return nsxModel.VpcSubnet{
		Id:          &subnetID,
		DisplayName: &subnetDisplayName,
		Description: &subnetDescription,
		Revision:    &subnetRevision,
	}
}

func minimalSubnetData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":     subnetDisplayName,
		"description":      subnetDescription,
		"nsx_id":           subnetID,
		"ip_blocks":        []interface{}{subnetIpBlock},
		"access_mode":      "Private",
		"ipv4_subnet_size": 16,
	}
}

func setupSubnetMock(t *testing.T, ctrl *gomock.Controller) (*vpcsmocks.MockSubnetsClient, *subnetsmocks.MockPortsClient, func()) {
	mockSubnetSDK := vpcsmocks.NewMockSubnetsClient(ctrl)
	mockPortsSDK := subnetsmocks.NewMockPortsClient(ctrl)

	mockSubnetWrapper := &apivpcs.VpcSubnetClientContext{
		Client:     mockSubnetSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}
	mockPortsWrapper := &apisubnets.VpcSubnetPortClientContext{
		Client:     mockPortsSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	origSubnet := cliVpcSubnetsClient
	origPorts := cliVpcSubnetPortsClient
	cliVpcSubnetsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcSubnetClientContext {
		return mockSubnetWrapper
	}
	cliVpcSubnetPortsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apisubnets.VpcSubnetPortClientContext {
		return mockPortsWrapper
	}
	return mockSubnetSDK, mockPortsSDK, func() {
		cliVpcSubnetsClient = origSubnet
		cliVpcSubnetPortsClient = origPorts
	}
}

func TestMockResourceNsxtVpcSubnetCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSubnetSDK, _, restore := setupSubnetMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSubnetSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), subnetID).Return(nsxModel.VpcSubnet{}, notFoundErr),
			mockSubnetSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), subnetID, gomock.Any()).Return(nil),
			mockSubnetSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), subnetID).Return(subnetAPIResponse(), nil),
		)

		res := resourceNsxtVpcSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSubnetData())

		err := resourceNsxtVpcSubnetCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, subnetID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSubnetData())

		err := resourceNsxtVpcSubnetCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcSubnetRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSubnetSDK, _, restore := setupSubnetMock(t, ctrl)
		defer restore()

		mockSubnetSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), subnetID).Return(subnetAPIResponse(), nil)

		res := resourceNsxtVpcSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSubnetData())
		d.SetId(subnetID)

		err := resourceNsxtVpcSubnetRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, subnetDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSubnetSDK, _, restore := setupSubnetMock(t, ctrl)
		defer restore()

		mockSubnetSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), subnetID).Return(nsxModel.VpcSubnet{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSubnetData())
		d.SetId(subnetID)

		err := resourceNsxtVpcSubnetRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcSubnetUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSubnetSDK, _, restore := setupSubnetMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSubnetSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), subnetID, gomock.Any()).Return(subnetAPIResponse(), nil),
			mockSubnetSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), subnetID).Return(subnetAPIResponse(), nil),
		)

		res := resourceNsxtVpcSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSubnetData())
		d.SetId(subnetID)

		err := resourceNsxtVpcSubnetUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcSubnetDelete(t *testing.T) {
	t.Run("Delete success with no ports", func(t *testing.T) {
		// Override poll delay for fast test
		vpcSubnetDeletePortsPollDelay = 0
		vpcSubnetDeletePortsPollMinTimeout = 0
		defer func() {
			vpcSubnetDeletePortsPollDelay = 1 * time.Second
			vpcSubnetDeletePortsPollMinTimeout = 1 * time.Second
		}()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSubnetSDK, mockPortsSDK, restore := setupSubnetMock(t, ctrl)
		defer restore()

		emptyPorts := nsxModel.VpcSubnetPortListResult{Results: []nsxModel.VpcSubnetPort{}}
		mockPortsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), subnetID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(emptyPorts, nil)
		mockSubnetSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), subnetID).Return(nil)

		res := resourceNsxtVpcSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSubnetData())
		d.SetId(subnetID)

		err := resourceNsxtVpcSubnetDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when port list returns error", func(t *testing.T) {
		vpcSubnetDeletePortsPollDelay = 0
		vpcSubnetDeletePortsPollMinTimeout = 0
		defer func() {
			vpcSubnetDeletePortsPollDelay = 1 * time.Second
			vpcSubnetDeletePortsPollMinTimeout = 1 * time.Second
		}()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockPortsSDK, restore := setupSubnetMock(t, ctrl)
		defer restore()

		mockPortsSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), subnetID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.VpcSubnetPortListResult{}, errors.New("list error"))

		res := resourceNsxtVpcSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSubnetData())
		d.SetId(subnetID)

		err := resourceNsxtVpcSubnetDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestValidateSubnetDhcpv6Config(t *testing.T) {
	res := resourceNsxtVpcSubnet()
	t.Run("DHCP_RELAY with dhcpv6_server_additional_config is rejected", func(t *testing.T) {
		data := minimalSubnetData()
		data["subnet_dhcpv6_config"] = []interface{}{map[string]interface{}{
			"mode": nsxModel.SubnetDhcpv6Config_MODE_RELAY,
			"dhcpv6_server_additional_config": []interface{}{map[string]interface{}{
				"domain_names": []interface{}{"example.com"},
			}},
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		err := validateSubnetDhcpv6Config(d)
		require.Error(t, err)
	})
	t.Run("DHCP_DEACTIVATED with reserved_ip_ranges is rejected", func(t *testing.T) {
		data := minimalSubnetData()
		data["subnet_dhcpv6_config"] = []interface{}{map[string]interface{}{
			"mode": nsxModel.SubnetDhcpv6Config_MODE_DEACTIVATED,
			"dhcpv6_server_additional_config": []interface{}{map[string]interface{}{
				"reserved_ip_ranges": []interface{}{"2001:db8::10"},
			}},
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		err := validateSubnetDhcpv6Config(d)
		require.Error(t, err)
	})
	t.Run("DHCP_SERVER with additional config is allowed", func(t *testing.T) {
		data := minimalSubnetData()
		data["subnet_dhcpv6_config"] = []interface{}{map[string]interface{}{
			"mode": nsxModel.SubnetDhcpv6Config_MODE_SERVER,
			"dhcpv6_server_additional_config": []interface{}{map[string]interface{}{
				"domain_names": []interface{}{"example.com"},
			}},
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		err := validateSubnetDhcpv6Config(d)
		require.NoError(t, err)
	})
}
