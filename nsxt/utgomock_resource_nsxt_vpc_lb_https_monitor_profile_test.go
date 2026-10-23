//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import "testing"

func TestMockResourceNsxtVpcLBHttpsMonitorProfile(t *testing.T) {
	runVpcLbProfileResourceTests(t, newVpcLbMonitorProfileResourceSpec(
		"nsxt_vpc_lb_https_monitor_profile",
		resourceNsxtVpcLBHttpsMonitorProfile,
		minimalLBHttpsMonitorData,
		lbHttpsMonitorStructValue,
		assertLBHttpsMonitorFields,
		lbHttpsMonitorID,
	))
}
