//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestMockDataSourceNsxtVpcLBMonitor(t *testing.T) {
	runVpcLbProfileDataSourceTests(t, vpcLbProfileDataSourceSpec{
		newDataSource: dataSourceNsxtVpcLBMonitor,
		setupMock:     setupVpcLbMonitorProfileMock,
		structValue:   lbMonitorProfileStructValue,
		resourceType:  nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBTCPMONITORPROFILE,
		typeName:      "TCP",
		otherTypeName: "HTTP",
		apiSegment:    "vpc-lb-monitor-profiles",
	})
}
