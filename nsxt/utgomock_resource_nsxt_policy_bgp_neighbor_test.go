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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	bgpapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services/bgp"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	bgpMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services/bgp"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	bgpNbID          = "bgp-nb-1"
	bgpNbBgpPath     = "/infra/tier-0s/t0-bgpnb/locale-services/default/bgp"
	bgpNbT0ID        = "t0-bgpnb"
	bgpNbServiceID   = "default"
	bgpNbAddr        = "192.168.1.1"
	bgpNbRemoteAsNum = "65000"
	bgpNbHoldDown    = int64(180)
	bgpNbKeepAlive   = int64(60)
	bgpNbMaxHopLimit = int64(255)
	bgpNbRevision    = int64(1)
	bgpNbPath        = "/infra/tier-0s/t0-bgpnb/locale-services/default/bgp/neighbors/bgp-nb-1"
	bgpNbName        = "bgp-nb-fooname"
)

func TestMockResourceNsxtPolicyBgpNeighborCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNeighborsSDK := bgpMocks.NewMockNeighborsClient(ctrl)
	nbWrapper := &bgpapi.BgpNeighborConfigClientContext{
		Client:     mockNeighborsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliBgpNeighborsClient
	defer func() { cliBgpNeighborsClient = originalCli }()
	cliBgpNeighborsClient = func(sessionContext utl.SessionContext, connector client.Connector) *bgpapi.BgpNeighborConfigClientContext {
		return nbWrapper
	}

	res := resourceNsxtPolicyBgpNeighbor()

	t.Run("Create_success", func(t *testing.T) {
		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, gomock.Any()).Return(model.BgpNeighborConfig{}, vapiErrors.NotFound{})
		mockNeighborsSDK.EXPECT().Patch(bgpNbT0ID, bgpNbServiceID, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, gomock.Any()).Return(model.BgpNeighborConfig{
			DisplayName:     &bgpNbName,
			Path:            &bgpNbPath,
			Revision:        &bgpNbRevision,
			NeighborAddress: &bgpNbAddr,
			RemoteAsNum:     &bgpNbRemoteAsNum,
			HoldDownTime:    &bgpNbHoldDown,
			KeepAliveTime:   &bgpNbKeepAlive,
			MaximumHopLimit: &bgpNbMaxHopLimit,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":     bgpNbName,
			"bgp_path":         bgpNbBgpPath,
			"neighbor_address": bgpNbAddr,
			"remote_as_num":    bgpNbRemoteAsNum,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, bgpNbName, d.Get("display_name"))
	})

	t.Run("Create_fails_when_already_exists", func(t *testing.T) {
		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, bgpNbID).Return(model.BgpNeighborConfig{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":           bgpNbID,
			"bgp_path":         bgpNbBgpPath,
			"neighbor_address": bgpNbAddr,
			"remote_as_num":    bgpNbRemoteAsNum,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_bgp_path_invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path":         "/invalid/path",
			"neighbor_address": bgpNbAddr,
			"remote_as_num":    bgpNbRemoteAsNum,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid bgp_path")
	})

	t.Run("Create_fails_with_source_attachment_before_NSX_9_2", func(t *testing.T) {
		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, gomock.Any()).Return(model.BgpNeighborConfig{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":      bgpNbName,
			"bgp_path":          bgpNbBgpPath,
			"neighbor_address":  bgpNbAddr,
			"remote_as_num":     bgpNbRemoteAsNum,
			"source_attachment": []interface{}{"/orgs/default/projects/p1/vpcs/v1/attachments/att-1"},
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, gomock.Any()).Return(model.BgpNeighborConfig{}, vapiErrors.NotFound{})
		mockNeighborsSDK.EXPECT().Patch(bgpNbT0ID, bgpNbServiceID, gomock.Any(), gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path":         bgpNbBgpPath,
			"neighbor_address": bgpNbAddr,
			"remote_as_num":    bgpNbRemoteAsNum,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyBgpNeighborRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNeighborsSDK := bgpMocks.NewMockNeighborsClient(ctrl)
	nbWrapper := &bgpapi.BgpNeighborConfigClientContext{
		Client:     mockNeighborsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliBgpNeighborsClient
	defer func() { cliBgpNeighborsClient = originalCli }()
	cliBgpNeighborsClient = func(sessionContext utl.SessionContext, connector client.Connector) *bgpapi.BgpNeighborConfigClientContext {
		return nbWrapper
	}

	res := resourceNsxtPolicyBgpNeighbor()

	t.Run("Read_success", func(t *testing.T) {
		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, bgpNbID).Return(model.BgpNeighborConfig{
			DisplayName:     &bgpNbName,
			Path:            &bgpNbPath,
			Revision:        &bgpNbRevision,
			NeighborAddress: &bgpNbAddr,
			RemoteAsNum:     &bgpNbRemoteAsNum,
			HoldDownTime:    &bgpNbHoldDown,
			KeepAliveTime:   &bgpNbKeepAlive,
			MaximumHopLimit: &bgpNbMaxHopLimit,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path": bgpNbBgpPath,
		})
		d.SetId(bgpNbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, bgpNbName, d.Get("display_name"))
		assert.Equal(t, bgpNbAddr, d.Get("neighbor_address"))
	})

	t.Run("Read_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path": bgpNbBgpPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BgpNeighbor ID")
	})

	t.Run("Read_fails_when_Get_returns_not_found", func(t *testing.T) {
		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, bgpNbID).Return(model.BgpNeighborConfig{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path": bgpNbBgpPath,
		})
		d.SetId(bgpNbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyBgpNeighborUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNeighborsSDK := bgpMocks.NewMockNeighborsClient(ctrl)
	nbWrapper := &bgpapi.BgpNeighborConfigClientContext{
		Client:     mockNeighborsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliBgpNeighborsClient
	defer func() { cliBgpNeighborsClient = originalCli }()
	cliBgpNeighborsClient = func(sessionContext utl.SessionContext, connector client.Connector) *bgpapi.BgpNeighborConfigClientContext {
		return nbWrapper
	}

	res := resourceNsxtPolicyBgpNeighbor()

	t.Run("Update_success", func(t *testing.T) {
		mockNeighborsSDK.EXPECT().Patch(bgpNbT0ID, bgpNbServiceID, bgpNbID, gomock.Any(), gomock.Any()).Return(nil)
		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, bgpNbID).Return(model.BgpNeighborConfig{
			DisplayName:     &bgpNbName,
			Path:            &bgpNbPath,
			Revision:        &bgpNbRevision,
			NeighborAddress: &bgpNbAddr,
			RemoteAsNum:     &bgpNbRemoteAsNum,
			HoldDownTime:    &bgpNbHoldDown,
			KeepAliveTime:   &bgpNbKeepAlive,
			MaximumHopLimit: &bgpNbMaxHopLimit,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path":         bgpNbBgpPath,
			"neighbor_address": bgpNbAddr,
			"remote_as_num":    bgpNbRemoteAsNum,
		})
		d.SetId(bgpNbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path": bgpNbBgpPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborUpdate(d, m)
		require.Error(t, err)
	})

	t.Run("Update_fails_when_Patch_returns_error", func(t *testing.T) {
		mockNeighborsSDK.EXPECT().Patch(bgpNbT0ID, bgpNbServiceID, bgpNbID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path":         bgpNbBgpPath,
			"neighbor_address": bgpNbAddr,
			"remote_as_num":    bgpNbRemoteAsNum,
		})
		d.SetId(bgpNbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyBgpNeighborDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNeighborsSDK := bgpMocks.NewMockNeighborsClient(ctrl)
	nbWrapper := &bgpapi.BgpNeighborConfigClientContext{
		Client:     mockNeighborsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliBgpNeighborsClient
	defer func() { cliBgpNeighborsClient = originalCli }()
	cliBgpNeighborsClient = func(sessionContext utl.SessionContext, connector client.Connector) *bgpapi.BgpNeighborConfigClientContext {
		return nbWrapper
	}

	res := resourceNsxtPolicyBgpNeighbor()

	t.Run("Delete_success", func(t *testing.T) {
		mockNeighborsSDK.EXPECT().Delete(bgpNbT0ID, bgpNbServiceID, bgpNbID, gomock.Any()).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path": bgpNbBgpPath,
		})
		d.SetId(bgpNbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path": bgpNbBgpPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborDelete(d, m)
		require.Error(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockNeighborsSDK.EXPECT().Delete(bgpNbT0ID, bgpNbServiceID, bgpNbID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path": bgpNbBgpPath,
		})
		d.SetId(bgpNbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborDelete(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyBgpNeighborRead_sourceAttachmentVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNeighborsSDK := bgpMocks.NewMockNeighborsClient(ctrl)
	nbWrapper := &bgpapi.BgpNeighborConfigClientContext{
		Client:     mockNeighborsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliBgpNeighborsClient
	defer func() { cliBgpNeighborsClient = originalCli }()
	cliBgpNeighborsClient = func(sessionContext utl.SessionContext, connector client.Connector) *bgpapi.BgpNeighborConfigClientContext {
		return nbWrapper
	}

	res := resourceNsxtPolicyBgpNeighbor()
	attPath := "/orgs/default/projects/p/vpcs/v/attachments/a1"

	t.Run("source_attachment omitted in state before NSX 9.2.0", func(t *testing.T) {
		util.NsxVersion = "3.2.0"
		defer func() { util.NsxVersion = "" }()

		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, bgpNbID).Return(model.BgpNeighborConfig{
			DisplayName:      &bgpNbName,
			Path:             &bgpNbPath,
			Revision:         &bgpNbRevision,
			NeighborAddress:  &bgpNbAddr,
			RemoteAsNum:      &bgpNbRemoteAsNum,
			HoldDownTime:     &bgpNbHoldDown,
			KeepAliveTime:    &bgpNbKeepAlive,
			MaximumHopLimit:  &bgpNbMaxHopLimit,
			SourceAttachment: []string{attPath},
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path": bgpNbBgpPath,
		})
		d.SetId(bgpNbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborRead(d, m)
		require.NoError(t, err)
		sa := d.Get("source_attachment")
		assert.Equal(t, 0, len(sa.([]interface{})))
	})

	t.Run("source_attachment set in state on NSX 9.2.0", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, bgpNbID).Return(model.BgpNeighborConfig{
			DisplayName:      &bgpNbName,
			Path:             &bgpNbPath,
			Revision:         &bgpNbRevision,
			NeighborAddress:  &bgpNbAddr,
			RemoteAsNum:      &bgpNbRemoteAsNum,
			HoldDownTime:     &bgpNbHoldDown,
			KeepAliveTime:    &bgpNbKeepAlive,
			MaximumHopLimit:  &bgpNbMaxHopLimit,
			SourceAttachment: []string{attPath},
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"bgp_path": bgpNbBgpPath,
		})
		d.SetId(bgpNbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpNeighborRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, attPath, d.Get("source_attachment").([]interface{})[0].(string))
	})
}

func TestMockResourceNsxtPolicyBgpNeighborCreate_sourceAttachment920(t *testing.T) {
	util.NsxVersion = "9.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNeighborsSDK := bgpMocks.NewMockNeighborsClient(ctrl)
	nbWrapper := &bgpapi.BgpNeighborConfigClientContext{
		Client:     mockNeighborsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliBgpNeighborsClient
	defer func() { cliBgpNeighborsClient = originalCli }()
	cliBgpNeighborsClient = func(sessionContext utl.SessionContext, connector client.Connector) *bgpapi.BgpNeighborConfigClientContext {
		return nbWrapper
	}

	res := resourceNsxtPolicyBgpNeighbor()
	att := "/orgs/default/projects/p1/vpcs/v1/attachments/att-1"

	mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, gomock.Any()).Return(model.BgpNeighborConfig{}, vapiErrors.NotFound{})
	mockNeighborsSDK.EXPECT().Patch(bgpNbT0ID, bgpNbServiceID, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, gomock.Any()).Return(model.BgpNeighborConfig{
		DisplayName:      &bgpNbName,
		Path:             &bgpNbPath,
		Revision:         &bgpNbRevision,
		NeighborAddress:  &bgpNbAddr,
		RemoteAsNum:      &bgpNbRemoteAsNum,
		HoldDownTime:     &bgpNbHoldDown,
		KeepAliveTime:    &bgpNbKeepAlive,
		MaximumHopLimit:  &bgpNbMaxHopLimit,
		SourceAttachment: []string{att},
	}, nil)

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"display_name":      bgpNbName,
		"bgp_path":          bgpNbBgpPath,
		"neighbor_address":  bgpNbAddr,
		"remote_as_num":     bgpNbRemoteAsNum,
		"source_attachment": []interface{}{att},
	})
	m := newGoMockProviderClient()
	err := resourceNsxtPolicyBgpNeighborCreate(d, m)
	require.NoError(t, err)
}

func TestUnitNsxt_resourceNsxtPolicyBgpNeighborImport(t *testing.T) {
	t.Run("wrong segment count is rejected", func(t *testing.T) {
		res := resourceNsxtPolicyBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("t0-1/default")

		_, err := resourceNsxtPolicyBgpNeighborImport(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tier0-id")
	})

	t.Run("valid path succeeds and sets bgp_path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockNeighborsSDK := bgpMocks.NewMockNeighborsClient(ctrl)
		nbWrapper := &bgpapi.BgpNeighborConfigClientContext{Client: mockNeighborsSDK, ClientType: utl.Local}
		originalCli := cliBgpNeighborsClient
		defer func() { cliBgpNeighborsClient = originalCli }()
		cliBgpNeighborsClient = func(sessionContext utl.SessionContext, connector client.Connector) *bgpapi.BgpNeighborConfigClientContext {
			return nbWrapper
		}
		parentPath := bgpNbBgpPath
		mockNeighborsSDK.EXPECT().Get(bgpNbT0ID, bgpNbServiceID, bgpNbID).Return(model.BgpNeighborConfig{ParentPath: &parentPath}, nil)

		res := resourceNsxtPolicyBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(bgpNbT0ID + "/" + bgpNbServiceID + "/" + bgpNbID)

		out, err := resourceNsxtPolicyBgpNeighborImport(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, bgpNbID, d.Id())
		assert.Equal(t, bgpNbBgpPath, d.Get("bgp_path"))
	})
}

func bgpNeighborMinimalData() map[string]interface{} {
	return map[string]interface{}{
		"neighbor_address": "10.0.0.1",
		"remote_as_num":    "65000",
	}
}

func TestUnitNsxt_resourceNsxtPolicyBgpNeighborResourceDataToStruct(t *testing.T) {
	res := resourceNsxtPolicyBgpNeighbor()

	t.Run("basic fields are converted", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, bgpNeighborMinimalData())
		obj, err := resourceNsxtPolicyBgpNeighborResourceDataToStruct(d, "neighbor-1")
		require.NoError(t, err)
		assert.Equal(t, "neighbor-1", *obj.Id)
		assert.Equal(t, "10.0.0.1", *obj.NeighborAddress)
		assert.Equal(t, "65000", *obj.RemoteAsNum)
	})

	t.Run("bfd_config, route_filtering, and local_as_config are converted", func(t *testing.T) {
		data := bgpNeighborMinimalData()
		data["bfd_config"] = []interface{}{
			map[string]interface{}{"enabled": true, "interval": 1000, "multiple": 5},
		}
		data["route_filtering"] = []interface{}{
			map[string]interface{}{
				"address_family":   "IPV4",
				"enabled":          true,
				"in_route_filter":  "/infra/tier-0s/t0/prefix-lists/pl1",
				"out_route_filter": "/infra/tier-0s/t0/prefix-lists/pl2",
				"maximum_routes":   100,
			},
		}
		data["neighbor_local_as_config"] = []interface{}{
			map[string]interface{}{"local_as_num": "65001", "as_path_modifier_type": "NO_PREPEND_REPLACE_AS"},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		obj, err := resourceNsxtPolicyBgpNeighborResourceDataToStruct(d, "neighbor-1")
		require.NoError(t, err)

		require.NotNil(t, obj.Bfd)
		assert.True(t, *obj.Bfd.Enabled)
		assert.EqualValues(t, 1000, *obj.Bfd.Interval)
		assert.EqualValues(t, 5, *obj.Bfd.Multiple)

		require.Len(t, obj.RouteFiltering, 1)
		assert.Equal(t, "IPV4", *obj.RouteFiltering[0].AddressFamily)
		assert.Equal(t, []string{"/infra/tier-0s/t0/prefix-lists/pl1"}, obj.RouteFiltering[0].InRouteFilters)
		assert.Equal(t, []string{"/infra/tier-0s/t0/prefix-lists/pl2"}, obj.RouteFiltering[0].OutRouteFilters)
		require.NotNil(t, obj.RouteFiltering[0].MaximumRoutes)
		assert.EqualValues(t, 100, *obj.RouteFiltering[0].MaximumRoutes)

		require.NotNil(t, obj.NeighborLocalAsConfig)
		assert.Equal(t, "65001", *obj.NeighborLocalAsConfig.LocalAsNum)
		assert.Equal(t, "NO_PREPEND_REPLACE_AS", *obj.NeighborLocalAsConfig.AsPathModifierType)
	})

	t.Run("source_attachment requires NSX 9.2.0 or higher", func(t *testing.T) {
		data := bgpNeighborMinimalData()
		data["source_attachment"] = []interface{}{"/infra/tier-0s/t0/tg-attachments/a1"}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		_, err := resourceNsxtPolicyBgpNeighborResourceDataToStruct(d, "neighbor-1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})

	t.Run("source_attachment is included when NSX version supports it", func(t *testing.T) {
		data := bgpNeighborMinimalData()
		data["source_attachment"] = []interface{}{"/infra/tier-0s/t0/tg-attachments/a1"}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		obj, err := resourceNsxtPolicyBgpNeighborResourceDataToStruct(d, "neighbor-1")
		require.NoError(t, err)
		assert.Equal(t, []string{"/infra/tier-0s/t0/tg-attachments/a1"}, obj.SourceAttachment)
	})
}

func TestUnitNsxt_setBgpNeighborConfigInSchema(t *testing.T) {
	res := resourceNsxtPolicyBgpNeighbor()

	t.Run("populates schema including nested bfd/route_filtering/local_as_config", func(t *testing.T) {
		d := res.TestResourceData()
		holdDown := int64(180)
		keepAlive := int64(60)
		maxHop := int64(1)
		neighborAddr := "10.0.0.1"
		remoteAs := "65000"
		gracefulMode := model.BgpNeighborConfig_GRACEFUL_RESTART_MODE_HELPER_ONLY
		bfdEnabled := true
		bfdInterval := int64(1000)
		bfdMultiple := int64(5)
		addrFamily := "IPV4"
		filterEnabled := true
		maxRoutes := int64(100)
		localAsNum := "65001"

		obj := model.BgpNeighborConfig{
			HoldDownTime:        &holdDown,
			KeepAliveTime:       &keepAlive,
			MaximumHopLimit:     &maxHop,
			NeighborAddress:     &neighborAddr,
			RemoteAsNum:         &remoteAs,
			GracefulRestartMode: &gracefulMode,
			Bfd: &model.BgpBfdConfig{
				Enabled:  &bfdEnabled,
				Interval: &bfdInterval,
				Multiple: &bfdMultiple,
			},
			RouteFiltering: []model.BgpRouteFiltering{
				{
					AddressFamily:   &addrFamily,
					Enabled:         &filterEnabled,
					InRouteFilters:  []string{"/infra/tier-0s/t0/prefix-lists/pl1"},
					OutRouteFilters: []string{"/infra/tier-0s/t0/prefix-lists/pl2"},
					MaximumRoutes:   &maxRoutes,
				},
			},
			NeighborLocalAsConfig: &model.BgpNeighborLocalAsConfig{
				LocalAsNum: &localAsNum,
			},
		}

		setBgpNeighborConfigInSchema(d, obj, false)

		assert.Equal(t, neighborAddr, d.Get("neighbor_address"))
		assert.Equal(t, 180, d.Get("hold_down_time"))

		bfd := d.Get("bfd_config").([]interface{})
		require.Len(t, bfd, 1)
		bfdMap := bfd[0].(map[string]interface{})
		assert.Equal(t, true, bfdMap["enabled"])
		assert.Equal(t, 1000, bfdMap["interval"])

		filters := d.Get("route_filtering").([]interface{})
		require.Len(t, filters, 1)
		filterMap := filters[0].(map[string]interface{})
		assert.Equal(t, "/infra/tier-0s/t0/prefix-lists/pl1", filterMap["in_route_filter"])
		assert.Equal(t, "/infra/tier-0s/t0/prefix-lists/pl2", filterMap["out_route_filter"])
		assert.Equal(t, 100, filterMap["maximum_routes"])

		localAsConfigs := d.Get("neighbor_local_as_config").([]interface{})
		require.Len(t, localAsConfigs, 1)
		assert.Equal(t, localAsNum, localAsConfigs[0].(map[string]interface{})["local_as_num"])
	})

	t.Run("nil bfd and local_as_config leave those schema fields empty", func(t *testing.T) {
		d := res.TestResourceData()
		holdDown := int64(180)
		keepAlive := int64(60)
		maxHop := int64(1)
		neighborAddr := "10.0.0.1"
		remoteAs := "65000"

		obj := model.BgpNeighborConfig{
			HoldDownTime:    &holdDown,
			KeepAliveTime:   &keepAlive,
			MaximumHopLimit: &maxHop,
			NeighborAddress: &neighborAddr,
			RemoteAsNum:     &remoteAs,
		}

		setBgpNeighborConfigInSchema(d, obj, false)

		assert.Empty(t, d.Get("bfd_config").([]interface{}))
		assert.Empty(t, d.Get("neighbor_local_as_config").([]interface{}))
	})
}
