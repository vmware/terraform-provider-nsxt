//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gmModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	tier1sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t1lsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
)

func setupTier1DsLocaleServicesMock(t *testing.T, ctrl *gomock.Controller) (*t1lsmocks.MockLocaleServicesClient, func()) {
	t.Helper()
	mockLS := t1lsmocks.NewMockLocaleServicesClient(ctrl)
	wrapper := &tier1sapi.LocaleServicesClientContext{Client: mockLS, ClientType: utl.Local}
	orig := cliTier1LocaleServicesClient
	cliTier1LocaleServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier1sapi.LocaleServicesClientContext {
		return wrapper
	}
	return mockLS, func() { cliTier1LocaleServicesClient = orig }
}

func TestUnitNsxt_dataSourceNsxtPolicyTier1GatewayRead(t *testing.T) {
	rt := "Tier1"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("t1-ds-1"), DisplayName: str("t1-ds-name"), Path: str("/infra/tier-1s/t1-ds-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockLS, restore := setupTier1DsLocaleServicesMock(t, ctrl)
		defer restore()

		edgeClusterPath := "/infra/sites/default/enforcement-points/default/edge-clusters/cl1"
		mockLS.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nsxModel.LocaleServices{EdgeClusterPath: &edgeClusterPath}, nil)

		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "t1-ds-1"})

		err := dataSourceNsxtPolicyTier1GatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "t1-ds-1", d.Id())
		assert.Equal(t, edgeClusterPath, d.Get("edge_cluster_path"))
	})

	t.Run("invalid path yields extraction error", func(t *testing.T) {
		ds := dataSourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"path": "/infra/tier-1s/"})

		err := dataSourceNsxtPolicyTier1GatewayRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not extract ID from path")
	})

	t.Run("global manager skips edge cluster lookup", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "t1-ds-1"})

		m := newGoMockProviderClient()
		m.PolicyGlobalManager = true
		err := dataSourceNsxtPolicyTier1GatewayRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, "t1-ds-1", d.Id())
		assert.Equal(t, "", d.Get("edge_cluster_path"))
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "t1-ds-1"})

		err := dataSourceNsxtPolicyTier1GatewayRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
