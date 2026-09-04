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
	gmModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestUnitNsxt_dataSourceNsxtPolicyL2VpnServiceRead(t *testing.T) {
	l2VpnServiceDsID := "l2-vpn-service-1"

	t.Run("success", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str(l2VpnServiceDsID), DisplayName: str("l2vpn-name"), Path: str("/infra/l2vpn"), ResourceType: str("L2VPNService"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":           l2VpnServiceDsID,
			"gateway_path": "/infra/tier-0s/t0",
		})

		err := dataSourceNsxtPolicyL2VpnServiceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l2VpnServiceDsID, d.Id())
	})

	t.Run("error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search failed")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": l2VpnServiceDsID})

		err := dataSourceNsxtPolicyL2VpnServiceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
