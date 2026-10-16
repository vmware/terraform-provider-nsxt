//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import "testing"

func TestMockResourceNsxtVpcLBUdpMonitorProfile(t *testing.T) {
	runVpcLbProfileResourceTests(t, newVpcLbMonitorProfileResourceSpec(
		"nsxt_vpc_lb_udp_monitor_profile",
		resourceNsxtVpcLBUdpMonitorProfile,
		minimalLBUdpMonitorData,
		lbUdpMonitorStructValue,
		assertLBUdpMonitorFields,
		lbUdpMonitorID,
	))
}
