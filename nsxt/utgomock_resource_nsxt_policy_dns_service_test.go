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
	projectmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dnsSvcID          = "dns-service-id"
	dnsSvcDisplayName = "test-dns-service"
	dnsSvcDescription = "Test Policy DNS Service"
	dnsSvcRevision    = int64(1)
)

func dnsSvcAPIResponse() nsxModel.PolicyDnsService {
	listenerIP := "/orgs/default/projects/p1/ip-address-allocations/ip1"
	vnsCluster := "/infra/sites/default/vns-clusters/cluster1"
	return nsxModel.PolicyDnsService{
		Id:                   &dnsSvcID,
		DisplayName:          &dnsSvcDisplayName,
		Description:          &dnsSvcDescription,
		Revision:             &dnsSvcRevision,
		AllocatedListenerIps: []string{listenerIP},
		VnsClusters:          []string{vnsCluster},
	}
}

func minimalDnsSvcData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":           dnsSvcDisplayName,
		"description":            dnsSvcDescription,
		"nsx_id":                 dnsSvcID,
		"allocated_listener_ips": []interface{}{"/orgs/default/projects/p1/ip-address-allocations/ip1"},
		"vns_clusters":           []interface{}{"/infra/sites/default/vns-clusters/cluster1"},
	}
}

func setupDnsSvcMock(t *testing.T, ctrl *gomock.Controller) (*projectmocks.MockDnsServicesClient, func()) {
	mockSDK := projectmocks.NewMockDnsServicesClient(ctrl)
	mockWrapper := &apiprojects.PolicyDnsServiceClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "project1",
	}

	original := cliPolicyDnsServicesClient
	cliPolicyDnsServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.PolicyDnsServiceClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliPolicyDnsServicesClient = original }
}

func TestMockResourceNsxtPolicyDnsServiceCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsSvcMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsSvcID).Return(nsxModel.PolicyDnsService{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), dnsSvcID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsSvcID).Return(dnsSvcAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsSvcData())

		err := resourceNsxtPolicyDnsServiceCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsSvcID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyDnsService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsSvcData())

		err := resourceNsxtPolicyDnsServiceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})
}

func TestMockResourceNsxtPolicyDnsServiceRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsSvcMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsSvcID).Return(dnsSvcAPIResponse(), nil)

		res := resourceNsxtPolicyDnsService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsSvcData())
		d.SetId(dnsSvcID)

		err := resourceNsxtPolicyDnsServiceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsSvcDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsSvcMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsSvcID).Return(nsxModel.PolicyDnsService{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyDnsService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsSvcData())
		d.SetId(dnsSvcID)

		err := resourceNsxtPolicyDnsServiceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyDnsServiceUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsSvcMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), dnsSvcID, gomock.Any()).Return(dnsSvcAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsSvcID).Return(dnsSvcAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsSvcData())
		d.SetId(dnsSvcID)

		err := resourceNsxtPolicyDnsServiceUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyDnsServiceDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsSvcMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), dnsSvcID).Return(nil)

		res := resourceNsxtPolicyDnsService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsSvcData())
		d.SetId(dnsSvcID)

		err := resourceNsxtPolicyDnsServiceDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
