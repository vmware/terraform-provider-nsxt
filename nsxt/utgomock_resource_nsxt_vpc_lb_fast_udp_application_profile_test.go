//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import "testing"

func TestMockResourceNsxtVpcLBFastUdpApplicationProfile(t *testing.T) {
	runVpcLbProfileResourceTests(t, newVpcLbAppProfileResourceSpec(
		"nsxt_vpc_lb_fast_udp_application_profile",
		resourceNsxtVpcLBFastUdpApplicationProfile,
		minimalLBFastUdpData,
		lbFastUdpStructValue,
		assertLBFastUdpFields,
		lbFastUdpID,
	))
}
