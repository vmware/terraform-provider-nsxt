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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	staticroutesapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/static_routes"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	bfdmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/static_routes"
)

var (
	bfdPeerID          = "bfd-peer-1"
	bfdPeerDisplayName = "Test BFD Peer"
	bfdPeerDescription = "Test BFD peer"
	bfdPeerRevision    = int64(1)
	bfdPeerGwPath      = "/infra/tier-0s/t0-gw-1"
	bfdPeerGwID        = "t0-gw-1"
	bfdPeerAddress     = "192.168.10.1"
	bfdPeerProfilePath = "/infra/bfd-profiles/default"
	bfdPeerPath        = "/infra/tier-0s/t0-gw-1/static-routes/bfd-peers/bfd-peer-1"
)

func bfdPeerAPIResponse() nsxModel.StaticRouteBfdPeer {
	enabled := true
	return nsxModel.StaticRouteBfdPeer{
		Id:             &bfdPeerID,
		DisplayName:    &bfdPeerDisplayName,
		Description:    &bfdPeerDescription,
		Revision:       &bfdPeerRevision,
		Path:           &bfdPeerPath,
		PeerAddress:    &bfdPeerAddress,
		BfdProfilePath: &bfdPeerProfilePath,
		Enabled:        &enabled,
	}
}

func minimalBfdPeerData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":     bfdPeerDisplayName,
		"description":      bfdPeerDescription,
		"nsx_id":           bfdPeerID,
		"gateway_path":     bfdPeerGwPath,
		"peer_address":     bfdPeerAddress,
		"bfd_profile_path": bfdPeerProfilePath,
		"enabled":          true,
	}
}

func setupBfdPeerMock(t *testing.T, ctrl *gomock.Controller) (*bfdmocks.MockBfdPeersClient, func()) {
	mockSDK := bfdmocks.NewMockBfdPeersClient(ctrl)
	mockWrapper := &staticroutesapi.StaticRouteBfdPeerClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliBfdPeersClient
	cliBfdPeersClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *staticroutesapi.StaticRouteBfdPeerClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliBfdPeersClient = original }
}

func TestMockResourceNsxtPolicyStaticRouteBfdPeerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupBfdPeerMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(bfdPeerGwID, bfdPeerID).Return(nsxModel.StaticRouteBfdPeer{}, notFoundErr),
			mockSDK.EXPECT().Patch(bfdPeerGwID, bfdPeerID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(bfdPeerGwID, bfdPeerID).Return(bfdPeerAPIResponse(), nil),
		)

		res := resourceNsxtPolicyStaticRouteBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBfdPeerData())

		err := resourceNsxtPolicyStaticRouteBfdPeerCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, bfdPeerID, d.Id())
		assert.Equal(t, bfdPeerDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when not Tier0 gateway", func(t *testing.T) {
		t1Data := minimalBfdPeerData()
		t1Data["gateway_path"] = "/infra/tier-1s/t1-gw-1"

		res := resourceNsxtPolicyStaticRouteBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, t1Data)

		err := resourceNsxtPolicyStaticRouteBfdPeerCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0")
	})
}

func TestMockResourceNsxtPolicyStaticRouteBfdPeerRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupBfdPeerMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(bfdPeerGwID, bfdPeerID).Return(bfdPeerAPIResponse(), nil)

		res := resourceNsxtPolicyStaticRouteBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBfdPeerData())
		d.SetId(bfdPeerID)

		err := resourceNsxtPolicyStaticRouteBfdPeerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, bfdPeerDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(bfdPeerGwID, bfdPeerID).Return(nsxModel.StaticRouteBfdPeer{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyStaticRouteBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBfdPeerData())
		d.SetId(bfdPeerID)

		err := resourceNsxtPolicyStaticRouteBfdPeerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyStaticRouteBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBfdPeerData())

		err := resourceNsxtPolicyStaticRouteBfdPeerRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyStaticRouteBfdPeerUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupBfdPeerMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(bfdPeerGwID, bfdPeerID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(bfdPeerGwID, bfdPeerID).Return(bfdPeerAPIResponse(), nil),
		)

		res := resourceNsxtPolicyStaticRouteBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBfdPeerData())
		d.SetId(bfdPeerID)

		err := resourceNsxtPolicyStaticRouteBfdPeerUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyStaticRouteBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBfdPeerData())

		err := resourceNsxtPolicyStaticRouteBfdPeerUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyStaticRouteBfdPeerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupBfdPeerMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(bfdPeerGwID, bfdPeerID).Return(nil)

		res := resourceNsxtPolicyStaticRouteBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBfdPeerData())
		d.SetId(bfdPeerID)

		err := resourceNsxtPolicyStaticRouteBfdPeerDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyStaticRouteBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBfdPeerData())

		err := resourceNsxtPolicyStaticRouteBfdPeerDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
