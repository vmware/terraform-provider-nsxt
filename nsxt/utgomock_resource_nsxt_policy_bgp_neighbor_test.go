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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

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
