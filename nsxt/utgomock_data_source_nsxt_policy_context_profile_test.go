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

func TestUnitNsxt_dataSourceNsxtPolicyContextProfileRead(t *testing.T) {
	contextProfileDsID := "context-profile-1"

	t.Run("success", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str(contextProfileDsID), DisplayName: str("cp-name"), Path: str("/infra/cp"), ResourceType: str("PolicyContextProfile"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyContextProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": contextProfileDsID})

		err := dataSourceNsxtPolicyContextProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, contextProfileDsID, d.Id())
	})

	t.Run("error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search failed")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyContextProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": contextProfileDsID})

		err := dataSourceNsxtPolicyContextProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
