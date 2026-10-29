//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"
)

func TestMockResourceNsxtVpcLBHttpsMonitorProfile(t *testing.T) {
	runVpcLbProfileResourceTests(t, newVpcLbMonitorProfileResourceSpec(
		"nsxt_vpc_lb_https_monitor_profile",
		resourceNsxtVpcLBHttpsMonitorProfile,
		minimalLBHttpsMonitorData,
		lbHttpsMonitorStructValue,
		assertLBHttpsMonitorFields,
		lbHttpsMonitorID,
	))

	t.Run("Read clears removed server SSL binding", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := setupVpcLbMonitorProfileMock(ctrl)
		defer mock.restore()

		removed := model.LBHttpsMonitorProfile{
			ResourceType: model.LBMonitorProfile_RESOURCE_TYPE_LBHTTPSMONITORPROFILE,
		}
		value, errs := bindings.NewTypeConverter().ConvertToVapi(removed, model.LBHttpsMonitorProfileBindingType())
		require.Empty(t, errs)
		response, ok := value.(*data.StructValue)
		require.True(t, ok)
		gomock.InOrder(
			mock.get(lbHttpsMonitorID, lbHttpsMonitorStructValue(t), nil),
			mock.get(lbHttpsMonitorID, response, nil),
		)

		res := resourceNsxtVpcLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbWithContext(minimalLBHttpsMonitorData()))
		d.SetId(lbHttpsMonitorID)
		client := newGoMockProviderClient()
		require.NoError(t, res.Read(d, client))
		require.Len(t, d.Get("server_ssl"), 1)
		require.NoError(t, res.Read(d, client))
		assert.Empty(t, d.Get("server_ssl"))
		assert.Equal(t, lbHttpsMonitorID, d.Id())
	})
}
