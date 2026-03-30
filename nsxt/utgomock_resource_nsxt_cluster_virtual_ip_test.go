// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/cluster/ApiVirtualIpClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/cluster/ApiVirtualIpClient.go ApiVirtualIpClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/cluster"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"

	clustermocks "github.com/vmware/terraform-provider-nsxt/mocks/cluster"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vipIPv4Address = "192.168.10.100"
	vipIPv6Address = "2001:db8::1"
)

func setupVirtualIPMock(t *testing.T, ctrl *gomock.Controller) (*clustermocks.MockApiVirtualIpClient, func()) {
	mockVipSDK := clustermocks.NewMockApiVirtualIpClient(ctrl)

	originalCli := cliApiVirtualIpClient
	cliApiVirtualIpClient = func(_ vapiProtocolClient.Connector) cluster.ApiVirtualIpClient {
		return mockVipSDK
	}
	return mockVipSDK, func() { cliApiVirtualIpClient = originalCli }
}

func clusterVirtualIPData() map[string]interface{} {
	return map[string]interface{}{
		"force":        true,
		"ip_address":   vipIPv4Address,
		"ipv6_address": vipIPv6Address,
	}
}

func virtualIPProperties(ipv4, ipv6 string) nsxModel.ClusterVirtualIpProperties {
	return nsxModel.ClusterVirtualIpProperties{
		IpAddress:  &ipv4,
		Ip6Address: &ipv6,
	}
}

func TestMockResourceNsxtClusterVirtualIPCreate(t *testing.T) {
	t.Run("Create success with NSX >= 4.0.0", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		mockVipSDK.EXPECT().Setvirtualip(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(virtualIPProperties(vipIPv4Address, vipIPv6Address), nil)
		mockVipSDK.EXPECT().Get().
			Return(virtualIPProperties(vipIPv4Address, vipIPv6Address), nil)

		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, clusterVirtualIPData())

		err := resourceNsxtClusterVirtualIPCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, vipIPv4Address, d.Get("ip_address"))
		assert.Equal(t, vipIPv6Address, d.Get("ipv6_address"))
	})

	t.Run("Create success with NSX < 4.0.0 (IPv4 only)", func(t *testing.T) {
		util.NsxVersion = "3.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		ipv4Only := vipIPv4Address
		defaultIPv6 := DefaultIPv6VirtualAddress
		mockVipSDK.EXPECT().Setvirtualip((*string)(nil), (*string)(nil), &ipv4Only).
			Return(virtualIPProperties(vipIPv4Address, defaultIPv6), nil)
		mockVipSDK.EXPECT().Get().
			Return(virtualIPProperties(vipIPv4Address, defaultIPv6), nil)

		res := resourceNsxtClusterVirtualIP()
		data := clusterVirtualIPData()
		data["ipv6_address"] = DefaultIPv6VirtualAddress
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtClusterVirtualIPCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create fails when Setvirtualip returns error", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		mockVipSDK.EXPECT().Setvirtualip(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.ClusterVirtualIpProperties{}, errors.New("API error"))

		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, clusterVirtualIPData())

		err := resourceNsxtClusterVirtualIPCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtClusterVirtualIPRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		mockVipSDK.EXPECT().Get().
			Return(virtualIPProperties(vipIPv4Address, vipIPv6Address), nil)

		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("some-id")

		err := resourceNsxtClusterVirtualIPRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vipIPv4Address, d.Get("ip_address"))
		assert.Equal(t, vipIPv6Address, d.Get("ipv6_address"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtClusterVirtualIPRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining ClusterVirtualIP ID")
	})

	t.Run("Read fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		mockVipSDK.EXPECT().Get().
			Return(nsxModel.ClusterVirtualIpProperties{}, errors.New("API error"))

		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("some-id")

		err := resourceNsxtClusterVirtualIPRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtClusterVirtualIPUpdate(t *testing.T) {
	t.Run("Update success with NSX >= 4.0.0", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		mockVipSDK.EXPECT().Setvirtualip(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(virtualIPProperties(vipIPv4Address, vipIPv6Address), nil)
		mockVipSDK.EXPECT().Get().
			Return(virtualIPProperties(vipIPv4Address, vipIPv6Address), nil)

		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, clusterVirtualIPData())
		d.SetId("some-id")

		err := resourceNsxtClusterVirtualIPUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vipIPv4Address, d.Get("ip_address"))
	})

	t.Run("Update fails when Setvirtualip returns error", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		mockVipSDK.EXPECT().Setvirtualip(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.ClusterVirtualIpProperties{}, errors.New("update API error"))

		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, clusterVirtualIPData())
		d.SetId("some-id")

		err := resourceNsxtClusterVirtualIPUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update API error")
	})
}

func TestMockResourceNsxtClusterVirtualIPDelete(t *testing.T) {
	t.Run("Delete success with NSX >= 4.0.0 (clears both IPv4 and IPv6)", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		mockVipSDK.EXPECT().Clearvirtualip().
			Return(virtualIPProperties(DefaultIPv4VirtualAddress, vipIPv6Address), nil)
		mockVipSDK.EXPECT().Clearvirtualip6().
			Return(virtualIPProperties(DefaultIPv4VirtualAddress, DefaultIPv6VirtualAddress), nil)

		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("some-id")

		err := resourceNsxtClusterVirtualIPDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete success with NSX < 4.0.0 (IPv4 only)", func(t *testing.T) {
		util.NsxVersion = "3.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		mockVipSDK.EXPECT().Clearvirtualip().
			Return(virtualIPProperties(DefaultIPv4VirtualAddress, DefaultIPv6VirtualAddress), nil)

		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("some-id")

		err := resourceNsxtClusterVirtualIPDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtClusterVirtualIPDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining ClusterVirtualIP ID")
	})

	t.Run("Delete fails when Clearvirtualip returns error", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		mockVipSDK.EXPECT().Clearvirtualip().
			Return(nsxModel.ClusterVirtualIpProperties{}, errors.New("clear IPv4 error"))

		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("some-id")

		err := resourceNsxtClusterVirtualIPDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "clear IPv4 error")
	})

	t.Run("Delete fails when Clearvirtualip6 returns error", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVipSDK, restore := setupVirtualIPMock(t, ctrl)
		defer restore()

		mockVipSDK.EXPECT().Clearvirtualip().
			Return(virtualIPProperties(DefaultIPv4VirtualAddress, vipIPv6Address), nil)
		mockVipSDK.EXPECT().Clearvirtualip6().
			Return(nsxModel.ClusterVirtualIpProperties{}, errors.New("clear IPv6 error"))

		res := resourceNsxtClusterVirtualIP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("some-id")

		err := resourceNsxtClusterVirtualIPDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "clear IPv6 error")
	})
}
