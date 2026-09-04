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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestUnitNsxt_dataSourceNsxtUplinkHostSwitchProfileRead(t *testing.T) {
	rt := infra.HostSwitchProfiles_LIST_HOSTSWITCH_PROFILE_TYPE_POLICYUPLINKHOSTSWITCHPROFILE
	// PolicyUplinkHostSwitchProfileBindingType declares resource_type as a plain
	// (non-optional) string field, unlike the generic gm PolicyResource binding type
	// used by policyResourceToStructValue. Build the fixture with the concrete model
	// type/binding so the data source's own secondary ConvertToGolang call succeeds.
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(nsxModel.PolicyUplinkHostSwitchProfile{
		Id: str("uplink-1"), DisplayName: str("uplink-name"), Path: str("/infra/host-switch-profiles/uplink-1"), ResourceType: rt,
	}, nsxModel.PolicyUplinkHostSwitchProfileBindingType())
	require.Empty(t, errs)
	sv := val.(*data.StructValue)

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "uplink-1"})

		err := dataSourceNsxtUplinkHostSwitchProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "uplink-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtUplinkHostSwitchProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "uplink-name"})

		err := dataSourceNsxtUplinkHostSwitchProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "encountered an error while searching PolicyUplinkHostSwitchProfile")
	})
}
