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

func TestUnitNsxt_dataSourceNsxtPolicyIPSecVpnTunnelProfileRead(t *testing.T) {
	ipsecTunnelProfileDsID := "ipsec-tunnel-profile-1"

	t.Run("success", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str(ipsecTunnelProfileDsID), DisplayName: str("tp-name"), Path: str("/infra/tp"), ResourceType: str("IPSecVpnTunnelProfile"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyIPSecVpnTunnelProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": ipsecTunnelProfileDsID})

		err := dataSourceNsxtPolicyIPSecVpnTunnelProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecTunnelProfileDsID, d.Id())
	})

	t.Run("error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search failed")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyIPSecVpnTunnelProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": ipsecTunnelProfileDsID})

		err := dataSourceNsxtPolicyIPSecVpnTunnelProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}
