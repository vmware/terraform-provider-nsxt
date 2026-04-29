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

	apiprojects "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	cnamocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
)

var (
	cnaID          = "cna-test-id"
	cnaDisplayName = "test-cna"
	cnaDescription = "Test CNA"
	cnaRevision    = int64(1)
	cnaPath        = "/orgs/default/projects/project1/centralized-network-attachments/cna-test-id"
	cnaOrgID       = "default"
	cnaProjectID   = "project1"
	cnaSubnetPath  = "/orgs/default/projects/project1/vpcs/vpc1/subnets/sub1"
)

func cnaAPIResponse() nsxModel.CentralizedNetworkAttachment {
	return nsxModel.CentralizedNetworkAttachment{
		Id:          &cnaID,
		DisplayName: &cnaDisplayName,
		Description: &cnaDescription,
		Revision:    &cnaRevision,
		Path:        &cnaPath,
		SubnetPath:  &cnaSubnetPath,
	}
}

func minimalCNAData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": cnaDisplayName,
		"description":  cnaDescription,
		"nsx_id":       cnaID,
		"subnet_path":  cnaSubnetPath,
		"context": []interface{}{
			map[string]interface{}{
				"project_id":  cnaProjectID,
				"vpc_id":      "",
				"from_global": false,
			},
		},
	}
}

func setupCNAMock(t *testing.T, ctrl *gomock.Controller) (*cnamocks.MockCentralizedNetworkAttachmentsClient, func()) {
	mockSDK := cnamocks.NewMockCentralizedNetworkAttachmentsClient(ctrl)
	mockWrapper := &apiprojects.CentralizedNetworkAttachmentClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  cnaProjectID,
	}

	original := cliCNAClient
	cliCNAClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.CentralizedNetworkAttachmentClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliCNAClient = original }
}

func TestMockResourceNsxtPolicyCNACreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCNAMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(nsxModel.CentralizedNetworkAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(cnaAPIResponse(), nil),
		)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNACreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cnaID, d.Id())
		assert.Equal(t, cnaSubnetPath, d.Get("subnet_path"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(cnaAPIResponse(), nil)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNACreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when Patch returns error", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(nsxModel.CentralizedNetworkAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNACreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create with advertise_outbound_networks and interface_subnet", func(t *testing.T) {
		allowPrivate := true
		notFoundErr := vapiErrors.NotFound{}
		var captured nsxModel.CentralizedNetworkAttachment
		gomock.InOrder(
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(nsxModel.CentralizedNetworkAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).DoAndReturn(
				func(_, _, _ string, obj nsxModel.CentralizedNetworkAttachment) error {
					captured = obj
					return nil
				}),
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(cnaAPIResponse(), nil),
		)

		data := minimalCNAData()
		data["advertise_outbound_networks"] = []interface{}{
			map[string]interface{}{
				"allow_private":         allowPrivate,
				"allow_external_blocks": []interface{}{"10.0.0.0/8"},
			},
		}
		data["interface_subnet"] = []interface{}{
			map[string]interface{}{
				"interface_ip_address": []interface{}{"192.168.1.10"},
				"prefix_length":        24,
				"ha_vip_ip_address":    "192.168.1.100",
			},
		}

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyProjectCNACreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.NotNil(t, captured.AdvertiseOutboundNetworks)
		assert.True(t, *captured.AdvertiseOutboundNetworks.AllowPrivate)
		assert.Equal(t, []string{"10.0.0.0/8"}, captured.AdvertiseOutboundNetworks.AllowExternalBlocks)
		require.Len(t, captured.InterfaceSubnets, 1)
		assert.Equal(t, int64(24), *captured.InterfaceSubnets[0].PrefixLength)
		assert.Equal(t, "192.168.1.100", *captured.InterfaceSubnets[0].HaVipIpAddress)
	})
}

func TestMockResourceNsxtPolicyCNARead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCNAMock(t, ctrl)
	defer restore()

	t.Run("Read success sets all fields", func(t *testing.T) {
		allowPrivate := true
		prefixLen := int64(24)
		haVip := "192.168.1.100"
		apiResp := cnaAPIResponse()
		apiResp.AdvertiseOutboundNetworks = &nsxModel.AdvertiseOutboundNetworks{
			AllowPrivate:        &allowPrivate,
			AllowExternalBlocks: []string{"10.0.0.0/8"},
		}
		apiResp.InterfaceSubnets = []nsxModel.CentralizedNetworkAttachmentInterfaceSubnet{
			{
				InterfaceIpAddress: []string{"192.168.1.10"},
				PrefixLength:       &prefixLen,
				HaVipIpAddress:     &haVip,
			},
		}
		mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(apiResp, nil)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())
		d.SetId(cnaID)

		err := resourceNsxtPolicyProjectCNARead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cnaDisplayName, d.Get("display_name"))
		assert.Equal(t, cnaSubnetPath, d.Get("subnet_path"))

		aon := d.Get("advertise_outbound_networks").([]interface{})
		require.Len(t, aon, 1)
		aonMap := aon[0].(map[string]interface{})
		assert.True(t, aonMap["allow_private"].(bool))
		assert.Equal(t, []interface{}{"10.0.0.0/8"}, aonMap["allow_external_blocks"])

		subnets := d.Get("interface_subnet").([]interface{})
		require.Len(t, subnets, 1)
		subMap := subnets[0].(map[string]interface{})
		assert.Equal(t, 24, subMap["prefix_length"])
		assert.Equal(t, haVip, subMap["ha_vip_ip_address"])
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(nsxModel.CentralizedNetworkAttachment{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())
		d.SetId(cnaID)

		err := resourceNsxtPolicyProjectCNARead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNARead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyCNAUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCNAMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(cnaOrgID, cnaProjectID, cnaID, gomock.Any()).Return(cnaAPIResponse(), nil),
			mockSDK.EXPECT().Get(cnaOrgID, cnaProjectID, cnaID).Return(cnaAPIResponse(), nil),
		)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())
		d.SetId(cnaID)

		err := resourceNsxtPolicyProjectCNAUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNAUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyCNADelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCNAMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(cnaOrgID, cnaProjectID, cnaID).Return(nil)

		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())
		d.SetId(cnaID)

		err := resourceNsxtPolicyProjectCNADelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCNAData())

		err := resourceNsxtPolicyProjectCNADelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
