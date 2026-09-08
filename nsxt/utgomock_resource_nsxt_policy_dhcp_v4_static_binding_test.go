//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	infrasegsapi "github.com/vmware/terraform-provider-nsxt/api/infra/segments"
	t1segsapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/segments"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	segsDhcpMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/segments"
	t1segsDhcpMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/segments"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func setupDhcpV4SegmentMock(ctrl *gomock.Controller) (*segsDhcpMocks.MockDhcpStaticBindingConfigsClient, func()) {
	mockSDK := segsDhcpMocks.NewMockDhcpStaticBindingConfigsClient(ctrl)
	wrapper := &infrasegsapi.StructValueClientContext{Client: mockSDK, ClientType: utl.Local}
	orig := cliSegmentsDhcpStaticBindingConfigsClient
	cliSegmentsDhcpStaticBindingConfigsClient = func(_ utl.SessionContext, _ client.Connector) *infrasegsapi.StructValueClientContext {
		return wrapper
	}
	return mockSDK, func() { cliSegmentsDhcpStaticBindingConfigsClient = orig }
}

func setupDhcpV4T1SegmentMock(ctrl *gomock.Controller) (*t1segsDhcpMocks.MockDhcpStaticBindingConfigsClient, func()) {
	mockSDK := t1segsDhcpMocks.NewMockDhcpStaticBindingConfigsClient(ctrl)
	wrapper := &t1segsapi.StructValueClientContext{Client: mockSDK, ClientType: utl.Local}
	orig := cliT1SegmentsDhcpStaticBindingConfigsClient
	cliT1SegmentsDhcpStaticBindingConfigsClient = func(_ utl.SessionContext, _ client.Connector) *t1segsapi.StructValueClientContext {
		return wrapper
	}
	return mockSDK, func() { cliT1SegmentsDhcpStaticBindingConfigsClient = orig }
}

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

func TestUnitNsxt_getDhcpOptsFromSchema(t *testing.T) {
	res := resourceNsxtPolicyDhcpV4StaticBinding()

	t.Run("nil when neither option block is set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		assert.Nil(t, getDhcpOptsFromSchema(d))
	})

	t.Run("parses option_121 and generic options", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"dhcp_option_121": []interface{}{
				map[string]interface{}{"network": "10.0.0.0/24", "next_hop": "10.0.0.1"},
			},
			"dhcp_generic_option": []interface{}{
				map[string]interface{}{"code": 43, "values": []interface{}{"foo"}},
			},
		})
		opts := getDhcpOptsFromSchema(d)
		require.NotNil(t, opts)
		require.NotNil(t, opts.Option121)
		require.Len(t, opts.Others, 1)
	})
}

func TestMockNsxt_getPolicyDchpStaticBindingOnSegment(t *testing.T) {
	t.Run("infra segment success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpV4SegmentMock(ctrl)
		defer restore()
		sv := dhcpV4StructValue(t)
		mockSDK.EXPECT().Get(dhcpV4SegmentID, dhcpV4ID).Return(sv, nil)

		obj, err := getPolicyDchpStaticBindingOnSegment(utl.SessionContext{ClientType: utl.Local}, dhcpV4ID, dhcpV4SegmentPath, nil)
		require.NoError(t, err)
		assert.NotNil(t, obj)
	})

	t.Run("fixed tier-1 segment success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpV4T1SegmentMock(ctrl)
		defer restore()
		sv := dhcpV4StructValue(t)
		mockSDK.EXPECT().Get("gw-1", dhcpV4SegmentID, dhcpV4ID).Return(sv, nil)

		obj, err := getPolicyDchpStaticBindingOnSegment(utl.SessionContext{ClientType: utl.Local}, dhcpV4ID, "/infra/tier-1s/gw-1/segments/"+dhcpV4SegmentID, nil)
		require.NoError(t, err)
		assert.NotNil(t, obj)
	})

	t.Run("invalid segment path errors", func(t *testing.T) {
		_, err := getPolicyDchpStaticBindingOnSegment(utl.SessionContext{ClientType: utl.Local}, dhcpV4ID, "/infra/tier-1s/gw-1", nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid Segment Path")
	})
}

func TestMockNsxt_resourceNsxtPolicyDhcpStaticBindingExistsOnSegment(t *testing.T) {
	t.Run("Get succeeds means it exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpV4SegmentMock(ctrl)
		defer restore()
		sv := dhcpV4StructValue(t)
		mockSDK.EXPECT().Get(dhcpV4SegmentID, dhcpV4ID).Return(sv, nil)

		exists, err := resourceNsxtPolicyDhcpStaticBindingExistsOnSegment(utl.SessionContext{ClientType: utl.Local}, dhcpV4ID, dhcpV4SegmentPath, nil)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("NotFound means it does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpV4SegmentMock(ctrl)
		defer restore()
		mockSDK.EXPECT().Get(dhcpV4SegmentID, dhcpV4ID).Return(nil, vapiErrors.NotFound{})

		exists, err := resourceNsxtPolicyDhcpStaticBindingExistsOnSegment(utl.SessionContext{ClientType: utl.Local}, dhcpV4ID, dhcpV4SegmentPath, nil)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("other errors propagate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDhcpV4SegmentMock(ctrl)
		defer restore()
		mockSDK.EXPECT().Get(dhcpV4SegmentID, dhcpV4ID).Return(nil, vapiErrors.InternalServerError{})

		_, err := resourceNsxtPolicyDhcpStaticBindingExistsOnSegment(utl.SessionContext{ClientType: utl.Local}, dhcpV4ID, dhcpV4SegmentPath, nil)
		require.Error(t, err)
	})
}

func TestUnitNsxt_nsxtSegmentResourceImporter(t *testing.T) {
	res := resourceNsxtPolicyDhcpV4StaticBinding()

	t.Run("full policy path is parsed directly", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dhcpV4SegmentPath + "/dhcp-static-binding-configs/" + dhcpV4ID)

		out, err := nsxtSegmentResourceImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, dhcpV4SegmentPath, out[0].Get("segment_path"))
	})

	t.Run("legacy segmentID/bindingID format", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dhcpV4SegmentID + "/" + dhcpV4ID)

		out, err := nsxtSegmentResourceImporter(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "/infra/segments/"+dhcpV4SegmentID, out[0].Get("segment_path"))
		assert.Equal(t, dhcpV4ID, out[0].Id())
	})

	t.Run("legacy gatewayID/segmentID/bindingID format", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("gw-1/" + dhcpV4SegmentID + "/" + dhcpV4ID)

		out, err := nsxtSegmentResourceImporter(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "/infra/tier-1s/gw-1/segments/"+dhcpV4SegmentID, out[0].Get("segment_path"))
		assert.Equal(t, dhcpV4ID, out[0].Id())
	})

	t.Run("too short legacy format errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("just-one-segment")

		_, err := nsxtSegmentResourceImporter(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Import format")
	})
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
