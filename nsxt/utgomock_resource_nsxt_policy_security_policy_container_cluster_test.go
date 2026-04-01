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

	secpolapi "github.com/vmware/terraform-provider-nsxt/api/infra/domains/security_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	secpolmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains/security_policies"
)

var (
	ccSpanID              = "cc-span-1"
	ccSpanDisplayName     = "Test Container Cluster Span"
	ccSpanDescription     = "Test container cluster span"
	ccSpanRevision        = int64(1)
	ccSpanDomain          = "default"
	ccSpanPolicyID        = "sp-001"
	ccSpanPolicyPath      = "/infra/domains/default/security-policies/sp-001"
	ccSpanContainerClPath = "/infra/container-clusters/cc-1"
	ccSpanPath            = "/infra/domains/default/security-policies/sp-001/container-cluster-span/cc-span-1"
)

func ccSpanAPIResponse() nsxModel.SecurityPolicyContainerCluster {
	return nsxModel.SecurityPolicyContainerCluster{
		Id:                   &ccSpanID,
		DisplayName:          &ccSpanDisplayName,
		Description:          &ccSpanDescription,
		Revision:             &ccSpanRevision,
		ContainerClusterPath: &ccSpanContainerClPath,
		Path:                 &ccSpanPath,
	}
}

func minimalCCSpanData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":           ccSpanDisplayName,
		"description":            ccSpanDescription,
		"nsx_id":                 ccSpanID,
		"policy_path":            ccSpanPolicyPath,
		"container_cluster_path": ccSpanContainerClPath,
	}
}

func setupCCSpanMock(t *testing.T, ctrl *gomock.Controller) (*secpolmocks.MockContainerClusterSpanClient, func()) {
	mockSDK := secpolmocks.NewMockContainerClusterSpanClient(ctrl)
	mockWrapper := &secpolapi.ContainerClusterSpanClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliContainerClusterSpanClient
	cliContainerClusterSpanClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *secpolapi.ContainerClusterSpanClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliContainerClusterSpanClient = original }
}

func TestMockResourceNsxtPolicySecurityPolicyContainerClusterCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCCSpanMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(ccSpanDomain, ccSpanPolicyID, ccSpanID).Return(nsxModel.SecurityPolicyContainerCluster{}, notFoundErr),
			mockSDK.EXPECT().Patch(ccSpanDomain, ccSpanPolicyID, ccSpanID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ccSpanDomain, ccSpanPolicyID, ccSpanID).Return(ccSpanAPIResponse(), nil),
		)

		res := resourceNsxtPolicySecurityPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCCSpanData())

		err := resourceNsxtPolicySecurityPolicyContainerClusterCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ccSpanID, d.Id())
		assert.Equal(t, ccSpanDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(ccSpanDomain, ccSpanPolicyID, ccSpanID).Return(ccSpanAPIResponse(), nil)

		res := resourceNsxtPolicySecurityPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCCSpanData())

		err := resourceNsxtPolicySecurityPolicyContainerClusterCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicySecurityPolicyContainerClusterRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCCSpanMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ccSpanDomain, ccSpanPolicyID, ccSpanID).Return(ccSpanAPIResponse(), nil)

		res := resourceNsxtPolicySecurityPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCCSpanData())
		d.SetId(ccSpanID)

		err := resourceNsxtPolicySecurityPolicyContainerClusterRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ccSpanDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(ccSpanDomain, ccSpanPolicyID, ccSpanID).Return(nsxModel.SecurityPolicyContainerCluster{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicySecurityPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCCSpanData())
		d.SetId(ccSpanID)

		err := resourceNsxtPolicySecurityPolicyContainerClusterRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySecurityPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCCSpanData())

		err := resourceNsxtPolicySecurityPolicyContainerClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySecurityPolicyContainerClusterUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCCSpanMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(ccSpanDomain, ccSpanPolicyID, ccSpanID, gomock.Any()).Return(ccSpanAPIResponse(), nil),
			mockSDK.EXPECT().Get(ccSpanDomain, ccSpanPolicyID, ccSpanID).Return(ccSpanAPIResponse(), nil),
		)

		res := resourceNsxtPolicySecurityPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCCSpanData())
		d.SetId(ccSpanID)

		err := resourceNsxtPolicySecurityPolicyContainerClusterUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySecurityPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCCSpanData())

		err := resourceNsxtPolicySecurityPolicyContainerClusterUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySecurityPolicyContainerClusterDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCCSpanMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ccSpanDomain, ccSpanPolicyID, ccSpanID).Return(nil)

		res := resourceNsxtPolicySecurityPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCCSpanData())
		d.SetId(ccSpanID)

		err := resourceNsxtPolicySecurityPolicyContainerClusterDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySecurityPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCCSpanData())

		err := resourceNsxtPolicySecurityPolicyContainerClusterDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
