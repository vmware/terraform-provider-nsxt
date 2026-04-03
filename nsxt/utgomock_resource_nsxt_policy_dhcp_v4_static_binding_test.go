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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	infrasegsapi "github.com/vmware/terraform-provider-nsxt/api/infra/segments"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	segsDhcpMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/segments"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dhcpV4ID          = "dhcpv4-1"
	dhcpV4SegmentPath = "/infra/segments/segment-1"
	dhcpV4SegmentID   = "segment-1"
	dhcpV4MacAddress  = "00:11:22:33:44:55"
	dhcpV4IPAddress   = "192.168.1.10"
	dhcpV4Name        = "dhcpv4-fooname"
	dhcpV4Revision    = int64(1)
	dhcpV4LeaseTime   = int64(86400)
	dhcpV4Path        = "/infra/segments/segment-1/dhcp-static-binding-configs/dhcpv4-1"
)

func dhcpV4StructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	leaseTime := dhcpV4LeaseTime
	obj := model.DhcpV4StaticBindingConfig{
		DisplayName:  &dhcpV4Name,
		IpAddress:    &dhcpV4IPAddress,
		MacAddress:   &dhcpV4MacAddress,
		LeaseTime:    &leaseTime,
		Path:         &dhcpV4Path,
		Revision:     &dhcpV4Revision,
		ResourceType: "DhcpV4StaticBindingConfig",
	}
	val, errs := converter.ConvertToVapi(obj, model.DhcpV4StaticBindingConfigBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestMockResourceNsxtPolicyDhcpV4StaticBindingCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDhcpSDK := segsDhcpMocks.NewMockDhcpStaticBindingConfigsClient(ctrl)
	dhcpWrapper := &infrasegsapi.StructValueClientContext{
		Client:     mockDhcpSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentsDhcpStaticBindingConfigsClient
	defer func() { cliSegmentsDhcpStaticBindingConfigsClient = originalCli }()
	cliSegmentsDhcpStaticBindingConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *infrasegsapi.StructValueClientContext {
		return dhcpWrapper
	}

	res := resourceNsxtPolicyDhcpV4StaticBinding()

	t.Run("Create_success", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Patch(dhcpV4SegmentID, gomock.Any(), gomock.Any()).Return(nil)
		sv := dhcpV4StructValue(t)
		mockDhcpSDK.EXPECT().Get(dhcpV4SegmentID, gomock.Any()).Return(sv, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": dhcpV4Name,
			"segment_path": dhcpV4SegmentPath,
			"mac_address":  dhcpV4MacAddress,
			"ip_address":   dhcpV4IPAddress,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV4StaticBindingCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Patch(dhcpV4SegmentID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV4SegmentPath,
			"mac_address":  dhcpV4MacAddress,
			"ip_address":   dhcpV4IPAddress,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV4StaticBindingCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyDhcpV4StaticBindingRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDhcpSDK := segsDhcpMocks.NewMockDhcpStaticBindingConfigsClient(ctrl)
	dhcpWrapper := &infrasegsapi.StructValueClientContext{
		Client:     mockDhcpSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentsDhcpStaticBindingConfigsClient
	defer func() { cliSegmentsDhcpStaticBindingConfigsClient = originalCli }()
	cliSegmentsDhcpStaticBindingConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *infrasegsapi.StructValueClientContext {
		return dhcpWrapper
	}

	res := resourceNsxtPolicyDhcpV4StaticBinding()

	t.Run("Read_success", func(t *testing.T) {
		sv := dhcpV4StructValue(t)
		mockDhcpSDK.EXPECT().Get(dhcpV4SegmentID, dhcpV4ID).Return(sv, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV4SegmentPath,
		})
		d.SetId(dhcpV4ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV4StaticBindingRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, dhcpV4Name, d.Get("display_name"))
		assert.Equal(t, dhcpV4IPAddress, d.Get("ip_address"))
		assert.Equal(t, dhcpV4MacAddress, d.Get("mac_address"))
	})

	t.Run("Read_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV4SegmentPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV4StaticBindingRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "DhcpV4 Static Binding Config ID")
	})

	t.Run("Read_fails_when_Get_returns_not_found", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Get(dhcpV4SegmentID, dhcpV4ID).Return(nil, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV4SegmentPath,
		})
		d.SetId(dhcpV4ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV4StaticBindingRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyDhcpV4StaticBindingUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDhcpSDK := segsDhcpMocks.NewMockDhcpStaticBindingConfigsClient(ctrl)
	dhcpWrapper := &infrasegsapi.StructValueClientContext{
		Client:     mockDhcpSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentsDhcpStaticBindingConfigsClient
	defer func() { cliSegmentsDhcpStaticBindingConfigsClient = originalCli }()
	cliSegmentsDhcpStaticBindingConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *infrasegsapi.StructValueClientContext {
		return dhcpWrapper
	}

	res := resourceNsxtPolicyDhcpV4StaticBinding()

	t.Run("Update_success", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Patch(dhcpV4SegmentID, dhcpV4ID, gomock.Any()).Return(nil)
		sv := dhcpV4StructValue(t)
		mockDhcpSDK.EXPECT().Get(dhcpV4SegmentID, dhcpV4ID).Return(sv, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV4SegmentPath,
			"mac_address":  dhcpV4MacAddress,
			"ip_address":   dhcpV4IPAddress,
		})
		d.SetId(dhcpV4ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV4StaticBindingUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_Patch_returns_error", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Patch(dhcpV4SegmentID, dhcpV4ID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV4SegmentPath,
			"mac_address":  dhcpV4MacAddress,
			"ip_address":   dhcpV4IPAddress,
		})
		d.SetId(dhcpV4ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV4StaticBindingUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyDhcpV4StaticBindingDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDhcpSDK := segsDhcpMocks.NewMockDhcpStaticBindingConfigsClient(ctrl)
	dhcpWrapper := &infrasegsapi.StructValueClientContext{
		Client:     mockDhcpSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentsDhcpStaticBindingConfigsClient
	defer func() { cliSegmentsDhcpStaticBindingConfigsClient = originalCli }()
	cliSegmentsDhcpStaticBindingConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *infrasegsapi.StructValueClientContext {
		return dhcpWrapper
	}

	res := resourceNsxtPolicyDhcpV4StaticBinding()

	t.Run("Delete_success", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Delete(dhcpV4SegmentID, dhcpV4ID).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV4SegmentPath,
		})
		d.SetId(dhcpV4ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpStaticBindingDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV4SegmentPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpStaticBindingDelete(d, m)
		require.Error(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Delete(dhcpV4SegmentID, dhcpV4ID).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV4SegmentPath,
		})
		d.SetId(dhcpV4ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpStaticBindingDelete(d, m)
		require.Error(t, err)
	})
}
