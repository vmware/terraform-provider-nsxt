//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestMockDataSourceNsxtVpcLBAppProfile(t *testing.T) {
	runVpcLbProfileDataSourceTests(t, vpcLbProfileDataSourceSpec{
		newDataSource: dataSourceNsxtVpcLBAppProfile,
		setupMock:     setupVpcLbAppProfileMock,
		structValue:   lbAppProfileStructValue,
		resourceType:  nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE,
		typeName:      "TCP",
		otherTypeName: "HTTP",
		apiSegment:    "vpc-lb-app-profiles",
	})
}
