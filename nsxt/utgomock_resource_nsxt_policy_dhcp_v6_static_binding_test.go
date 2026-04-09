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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	infrasegsapi "github.com/vmware/terraform-provider-nsxt/api/infra/segments"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	segsDhcpMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/segments"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dhcpV6ID          = "dhcpv6-1"
	dhcpV6SegmentPath = "/infra/segments/segment-v6"
	dhcpV6SegmentID   = "segment-v6"
	dhcpV6MacAddress  = "00:11:22:33:44:66"
	dhcpV6Name        = "dhcpv6-fooname"
	dhcpV6Revision    = int64(1)
	dhcpV6LeaseTime   = int64(86400)
	dhcpV6Path        = "/infra/segments/segment-v6/dhcp-static-binding-configs/dhcpv6-1"
)

func dhcpV6StructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	leaseTime := dhcpV6LeaseTime
	obj := model.DhcpV6StaticBindingConfig{
		DisplayName:  &dhcpV6Name,
		MacAddress:   &dhcpV6MacAddress,
		LeaseTime:    &leaseTime,
		Path:         &dhcpV6Path,
		Revision:     &dhcpV6Revision,
		ResourceType: "DhcpV6StaticBindingConfig",
	}
	val, errs := converter.ConvertToVapi(obj, model.DhcpV6StaticBindingConfigBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestMockResourceNsxtPolicyDhcpV6StaticBindingCreate(t *testing.T) {
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

	res := resourceNsxtPolicyDhcpV6StaticBinding()

	t.Run("Create_success", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Patch(dhcpV6SegmentID, gomock.Any(), gomock.Any()).Return(nil)
		sv := dhcpV6StructValue(t)
		mockDhcpSDK.EXPECT().Get(dhcpV6SegmentID, gomock.Any()).Return(sv, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": dhcpV6Name,
			"segment_path": dhcpV6SegmentPath,
			"mac_address":  dhcpV6MacAddress,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV6StaticBindingCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Patch(dhcpV6SegmentID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV6SegmentPath,
			"mac_address":  dhcpV6MacAddress,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV6StaticBindingCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyDhcpV6StaticBindingRead(t *testing.T) {
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

	res := resourceNsxtPolicyDhcpV6StaticBinding()

	t.Run("Read_success", func(t *testing.T) {
		sv := dhcpV6StructValue(t)
		mockDhcpSDK.EXPECT().Get(dhcpV6SegmentID, dhcpV6ID).Return(sv, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV6SegmentPath,
		})
		d.SetId(dhcpV6ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV6StaticBindingRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, dhcpV6Name, d.Get("display_name"))
		assert.Equal(t, dhcpV6MacAddress, d.Get("mac_address"))
	})

	t.Run("Read_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV6SegmentPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV6StaticBindingRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "DhcpV6 Static Binding Config ID")
	})

	t.Run("Read_fails_when_Get_returns_not_found", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Get(dhcpV6SegmentID, dhcpV6ID).Return(nil, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV6SegmentPath,
		})
		d.SetId(dhcpV6ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV6StaticBindingRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyDhcpV6StaticBindingUpdate(t *testing.T) {
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

	res := resourceNsxtPolicyDhcpV6StaticBinding()

	t.Run("Update_success", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Patch(dhcpV6SegmentID, dhcpV6ID, gomock.Any()).Return(nil)
		sv := dhcpV6StructValue(t)
		mockDhcpSDK.EXPECT().Get(dhcpV6SegmentID, dhcpV6ID).Return(sv, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV6SegmentPath,
			"mac_address":  dhcpV6MacAddress,
		})
		d.SetId(dhcpV6ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV6StaticBindingUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_Patch_returns_error", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Patch(dhcpV6SegmentID, dhcpV6ID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV6SegmentPath,
			"mac_address":  dhcpV6MacAddress,
		})
		d.SetId(dhcpV6ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpV6StaticBindingUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyDhcpV6StaticBindingDelete(t *testing.T) {
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

	res := resourceNsxtPolicyDhcpV6StaticBinding()

	t.Run("Delete_success", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Delete(dhcpV6SegmentID, dhcpV6ID).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV6SegmentPath,
		})
		d.SetId(dhcpV6ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpStaticBindingDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Delete(dhcpV6SegmentID, dhcpV6ID).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": dhcpV6SegmentPath,
		})
		d.SetId(dhcpV6ID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpStaticBindingDelete(d, m)
		require.Error(t, err)
	})
}
