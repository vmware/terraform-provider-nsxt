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

	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tgwbfdmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
)

var (
	tgwBfdPeerID          = "bfd-peer-1"
	tgwBfdPeerDisplayName = "Test TGW BFD Peer"
	tgwBfdPeerDescription = "Test transit gateway BFD peer"
	tgwBfdPeerRevision    = int64(1)
	tgwBfdPeerParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1"
	tgwBfdPeerOrgID       = "default"
	tgwBfdPeerProjectID   = "project1"
	tgwBfdPeerTGWID       = "tgw1"
	tgwBfdPeerAddress     = "192.168.10.1"
	tgwBfdPeerProfilePath = "/orgs/default/projects/project1/infra/bfd-profiles/default"
	tgwBfdPeerPath        = "/orgs/default/projects/project1/transit-gateways/tgw1/bfd-peers/bfd-peer-1"
	tgwBfdPeerAttachment  = "/orgs/default/projects/project1/transit-gateways/tgw1/attachments/att1"
)

func tgwBfdPeerAPIResponse() nsxModel.TransitGatewayBfdPeer {
	enabled := true
	return nsxModel.TransitGatewayBfdPeer{
		Id:               &tgwBfdPeerID,
		DisplayName:      &tgwBfdPeerDisplayName,
		Description:      &tgwBfdPeerDescription,
		Revision:         &tgwBfdPeerRevision,
		Path:             &tgwBfdPeerPath,
		PeerAddress:      &tgwBfdPeerAddress,
		BfdProfilePath:   &tgwBfdPeerProfilePath,
		Enabled:          &enabled,
		SourceAttachment: []string{tgwBfdPeerAttachment},
	}
}

func minimalTGWBfdPeerData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":      tgwBfdPeerDisplayName,
		"description":       tgwBfdPeerDescription,
		"nsx_id":            tgwBfdPeerID,
		"parent_path":       tgwBfdPeerParentPath,
		"peer_address":      tgwBfdPeerAddress,
		"bfd_profile_path":  tgwBfdPeerProfilePath,
		"enabled":           true,
		"source_attachment": []interface{}{tgwBfdPeerAttachment},
	}
}

func setupTGWBfdPeerMock(t *testing.T, ctrl *gomock.Controller) (*tgwbfdmocks.MockBfdPeersClient, func()) {
	mockSDK := tgwbfdmocks.NewMockBfdPeersClient(ctrl)
	mockWrapper := &transitgateways.TransitGatewayBfdPeerClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwBfdPeerProjectID,
	}

	original := cliTGWBfdPeersClient
	cliTGWBfdPeersClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *transitgateways.TransitGatewayBfdPeerClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTGWBfdPeersClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayBfdPeerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWBfdPeerMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID).Return(nsxModel.TransitGatewayBfdPeer{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID).Return(tgwBfdPeerAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBfdPeerData())

		err := resourceNsxtPolicyTransitGatewayBfdPeerCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwBfdPeerID, d.Id())
		assert.Equal(t, tgwBfdPeerDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when parent_path is missing", func(t *testing.T) {
		data := minimalTGWBfdPeerData()
		data["parent_path"] = ""

		res := resourceNsxtPolicyTransitGatewayBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayBfdPeerCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayBfdPeerRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWBfdPeerMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID).Return(tgwBfdPeerAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGatewayBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBfdPeerData())
		d.SetId(tgwBfdPeerID)

		err := resourceNsxtPolicyTransitGatewayBfdPeerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwBfdPeerDisplayName, d.Get("display_name"))
		assert.Equal(t, tgwBfdPeerAddress, d.Get("peer_address"))
		assert.Equal(t, tgwBfdPeerProfilePath, d.Get("bfd_profile_path"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID).Return(nsxModel.TransitGatewayBfdPeer{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBfdPeerData())
		d.SetId(tgwBfdPeerID)

		err := resourceNsxtPolicyTransitGatewayBfdPeerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBfdPeerData())

		err := resourceNsxtPolicyTransitGatewayBfdPeerRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayBfdPeerUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWBfdPeerMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID, gomock.Any()).Return(tgwBfdPeerAPIResponse(), nil),
			mockSDK.EXPECT().Get(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID).Return(tgwBfdPeerAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBfdPeerData())
		d.SetId(tgwBfdPeerID)

		err := resourceNsxtPolicyTransitGatewayBfdPeerUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBfdPeerData())

		err := resourceNsxtPolicyTransitGatewayBfdPeerUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayBfdPeerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWBfdPeerMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBfdPeerData())
		d.SetId(tgwBfdPeerID)

		err := resourceNsxtPolicyTransitGatewayBfdPeerDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBfdPeerData())

		err := resourceNsxtPolicyTransitGatewayBfdPeerDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayBfdPeerWithoutOptional(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWBfdPeerMock(t, ctrl)
	defer restore()

	t.Run("Create without bfd_profile_path succeeds", func(t *testing.T) {
		data := map[string]interface{}{
			"display_name":      tgwBfdPeerDisplayName,
			"description":       tgwBfdPeerDescription,
			"nsx_id":            tgwBfdPeerID,
			"parent_path":       tgwBfdPeerParentPath,
			"peer_address":      tgwBfdPeerAddress,
			"enabled":           true,
			"source_attachment": []interface{}{tgwBfdPeerAttachment},
		}

		responseWithoutProfile := nsxModel.TransitGatewayBfdPeer{
			Id:               &tgwBfdPeerID,
			DisplayName:      &tgwBfdPeerDisplayName,
			Description:      &tgwBfdPeerDescription,
			Revision:         &tgwBfdPeerRevision,
			Path:             &tgwBfdPeerPath,
			PeerAddress:      &tgwBfdPeerAddress,
			SourceAttachment: []string{tgwBfdPeerAttachment},
		}

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID).Return(nsxModel.TransitGatewayBfdPeer{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwBfdPeerOrgID, tgwBfdPeerProjectID, tgwBfdPeerTGWID, tgwBfdPeerID).Return(responseWithoutProfile, nil),
		)

		res := resourceNsxtPolicyTransitGatewayBfdPeer()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayBfdPeerCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwBfdPeerID, d.Id())
		assert.Equal(t, "", d.Get("bfd_profile_path"))
	})
}
