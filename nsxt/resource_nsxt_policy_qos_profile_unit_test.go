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
)

func TestUnitNsxt_getPolicyQosRateShaperFromSchema(t *testing.T) {
	res := resourceNsxtPolicyQosProfile()

	t.Run("no shaper configured returns nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		sv := getPolicyQosRateShaperFromSchema(d, ingressRateShaperIndex)
		assert.Nil(t, sv)
	})

	t.Run("ingress rate shaper is converted to a StructValue", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"ingress_rate_shaper": []interface{}{
				map[string]interface{}{
					"enabled":         true,
					"average_bw_mbps": 10,
					"burst_size":      1000,
					"peak_bw_mbps":    20,
				},
			},
		})
		sv := getPolicyQosRateShaperFromSchema(d, ingressRateShaperIndex)
		require.NotNil(t, sv)

		converter := bindings.NewTypeConverter()
		golang, errs := converter.ConvertToGolang(sv, model.IngressRateLimiterBindingType())
		require.Empty(t, errs)
		shaper := golang.(model.IngressRateLimiter)
		assert.Equal(t, model.QosBaseRateLimiter_RESOURCE_TYPE_INGRESSRATELIMITER, shaper.ResourceType)
		assert.True(t, *shaper.Enabled)
		assert.EqualValues(t, 10, *shaper.AverageBandwidth)
		assert.EqualValues(t, 1000, *shaper.BurstSize)
		assert.EqualValues(t, 20, *shaper.PeakBandwidth)
	})

	t.Run("egress rate shaper uses the mbps scale and egress resource type", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"egress_rate_shaper": []interface{}{
				map[string]interface{}{
					"enabled":         false,
					"average_bw_mbps": 5,
					"burst_size":      500,
					"peak_bw_mbps":    15,
				},
			},
		})
		sv := getPolicyQosRateShaperFromSchema(d, egressRateShaperIndex)
		require.NotNil(t, sv)

		converter := bindings.NewTypeConverter()
		golang, errs := converter.ConvertToGolang(sv, model.IngressRateLimiterBindingType())
		require.Empty(t, errs)
		shaper := golang.(model.IngressRateLimiter)
		assert.Equal(t, model.QosBaseRateLimiter_RESOURCE_TYPE_EGRESSRATELIMITER, shaper.ResourceType)
		assert.False(t, *shaper.Enabled)
	})
}

func TestUnitNsxt_setPolicyQosRateShaperInSchema(t *testing.T) {
	res := resourceNsxtPolicyQosProfile()

	buildShaperStructValue := func(t *testing.T, resourceType string, enabled bool, avgBW, burstSize, peakBW int64) *data.StructValue {
		t.Helper()
		converter := bindings.NewTypeConverter()
		val, errs := converter.ConvertToVapi(model.IngressRateLimiter{
			ResourceType:     resourceType,
			Enabled:          &enabled,
			AverageBandwidth: &avgBW,
			BurstSize:        &burstSize,
			PeakBandwidth:    &peakBW,
		}, model.IngressRateLimiterBindingType())
		require.Empty(t, errs)
		return val.(*data.StructValue)
	}

	t.Run("matching resource type populates the schema", func(t *testing.T) {
		d := res.TestResourceData()
		sv := buildShaperStructValue(t, model.QosBaseRateLimiter_RESOURCE_TYPE_INGRESSRATELIMITER, true, 10, 1000, 20)

		setPolicyQosRateShaperInSchema(d, []*data.StructValue{sv}, ingressRateShaperIndex)

		shapers := d.Get("ingress_rate_shaper").([]interface{})
		require.Len(t, shapers, 1)
		elem := shapers[0].(map[string]interface{})
		assert.Equal(t, true, elem["enabled"])
		assert.EqualValues(t, 1000, elem["burst_size"])
		assert.EqualValues(t, 10, elem["average_bw_mbps"])
		assert.EqualValues(t, 20, elem["peak_bw_mbps"])
	})

	t.Run("mismatched resource type leaves the schema empty", func(t *testing.T) {
		d := res.TestResourceData()
		sv := buildShaperStructValue(t, model.QosBaseRateLimiter_RESOURCE_TYPE_EGRESSRATELIMITER, true, 10, 1000, 20)

		setPolicyQosRateShaperInSchema(d, []*data.StructValue{sv}, ingressRateShaperIndex)

		assert.Empty(t, d.Get("ingress_rate_shaper").([]interface{}))
	})
}
