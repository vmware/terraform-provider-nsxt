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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestUnitNsxt_dataSourceNsxtVtepHAHostSwitchProfileRead(t *testing.T) {
	rt := infra.HostSwitchProfiles_LIST_HOSTSWITCH_PROFILE_TYPE_POLICYVTEPHAHOSTSWITCHPROFILE
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("vtepha-1"), DisplayName: str("vtepha-name"), Path: str("/infra/host-switch-profiles/vtepha-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "vtepha-1"})

		err := dataSourceNsxtVtepHAHostSwitchProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "vtepha-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtVtepHAHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "vtepha-name"})

		err := dataSourceNsxtVtepHAHostSwitchProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}
