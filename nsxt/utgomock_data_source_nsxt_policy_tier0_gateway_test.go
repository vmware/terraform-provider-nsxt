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

	tier0localeservices "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0lsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
)

func setupTier0DsLocaleServicesMock(t *testing.T, ctrl *gomock.Controller) (*t0lsmocks.MockLocaleServicesClient, func()) {
	t.Helper()
	mockLS := t0lsmocks.NewMockLocaleServicesClient(ctrl)
	wrapper := &tier0localeservices.LocaleServicesClientContext{Client: mockLS, ClientType: utl.Local}
	orig := cliTier0LocaleServicesClient
	cliTier0LocaleServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0localeservices.LocaleServicesClientContext {
		return wrapper
	}
	return mockLS, func() { cliTier0LocaleServicesClient = orig }
}

func TestUnitNsxt_dataSourceNsxtPolicyTier0GatewayRead(t *testing.T) {
	rt := "Tier0"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("t0-ds-1"), DisplayName: str("t0-ds-name"), Path: str("/infra/tier-0s/t0-ds-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockLS, restore := setupTier0DsLocaleServicesMock(t, ctrl)
		defer restore()

		edgeClusterPath := "/infra/sites/default/enforcement-points/default/edge-clusters/cl1"
		mockLS.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nsxModel.LocaleServices{EdgeClusterPath: &edgeClusterPath}, nil)

		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "t0-ds-1"})

		err := dataSourceNsxtPolicyTier0GatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "t0-ds-1", d.Id())
		assert.Equal(t, edgeClusterPath, d.Get("edge_cluster_path"))
	})

	t.Run("global manager skips edge cluster lookup", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "t0-ds-1"})

		m := newGoMockProviderClient()
		m.PolicyGlobalManager = true
		err := dataSourceNsxtPolicyTier0GatewayRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, "t0-ds-1", d.Id())
		assert.Equal(t, "", d.Get("edge_cluster_path"))
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "t0-ds-1"})

		err := dataSourceNsxtPolicyTier0GatewayRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
