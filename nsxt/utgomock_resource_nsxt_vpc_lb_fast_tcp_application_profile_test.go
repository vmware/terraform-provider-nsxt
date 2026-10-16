//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import "testing"

func TestMockResourceNsxtVpcLBFastTcpApplicationProfile(t *testing.T) {
	runVpcLbProfileResourceTests(t, newVpcLbAppProfileResourceSpec(
		"nsxt_vpc_lb_fast_tcp_application_profile",
		resourceNsxtVpcLBFastTcpApplicationProfile,
		minimalLBFastTcpData,
		lbFastTcpStructValue,
		assertLBFastTcpFields,
		lbFastTcpID,
	))
}
