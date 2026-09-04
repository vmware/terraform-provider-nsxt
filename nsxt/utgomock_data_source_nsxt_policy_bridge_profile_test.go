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

func TestUnitNsxt_dataSourceNsxtPolicyBridgeProfileRead(t *testing.T) {
	bridgeProfileDsID := "bridge-profile-1"

	t.Run("success", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str(bridgeProfileDsID), DisplayName: str("bp-name"), Path: str("/infra/bp"), ResourceType: str("L2BridgeEndpointProfile"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyBridgeProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": bridgeProfileDsID})

		err := dataSourceNsxtPolicyBridgeProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, bridgeProfileDsID, d.Id())
	})

	t.Run("error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search failed")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyBridgeProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": bridgeProfileDsID})

		err := dataSourceNsxtPolicyBridgeProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
