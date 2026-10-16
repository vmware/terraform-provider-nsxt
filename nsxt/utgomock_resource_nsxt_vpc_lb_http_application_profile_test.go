//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import "testing"

func TestMockResourceNsxtVpcLBHttpApplicationProfile(t *testing.T) {
	runVpcLbProfileResourceTests(t, newVpcLbAppProfileResourceSpec(
		"nsxt_vpc_lb_http_application_profile",
		resourceNsxtVpcLBHttpApplicationProfile,
		minimalLBHttpAppData,
		lbHttpAppStructValue,
		assertLBHttpAppFields,
		lbHttpAppID,
	))
}
