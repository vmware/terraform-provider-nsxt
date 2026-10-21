//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import "testing"

func TestMockResourceNsxtVpcLBPassiveMonitorProfile(t *testing.T) {
	runVpcLbProfileResourceTests(t, newVpcLbMonitorProfileResourceSpec(
		"nsxt_vpc_lb_passive_monitor_profile",
		resourceNsxtVpcLBPassiveMonitorProfile,
		minimalLBPassiveMonitorData,
		lbPassiveMonitorStructValue,
		assertLBPassiveMonitorFields,
		lbPassiveMonitorID,
	))
}
