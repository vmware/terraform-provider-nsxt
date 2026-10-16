//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// The shared suite and helpers are defined in
// utgomock_resource_nsxt_vpc_lb_cookie_persistence_profile_test.go

func TestMockResourceNsxtVpcLBSourceIpPersistenceProfile(t *testing.T) {
	runVpcLbPersistenceResourceSuite(t, vpcLbPersistCase{
		name:        "source_ip",
		newRes:      resourceNsxtVpcLBSourceIpPersistenceProfile,
		id:          "vpc-lb-source-ip-1",
		displayName: "Test VPC LB Source IP Persistence Profile",
		extraData: map[string]interface{}{
			"purge":   "FULL",
			"timeout": 300,
		},
		newModel: func(displayName string) (interface{}, bindings.BindingType) {
			return model.LBSourceIpPersistenceProfile{
				DisplayName:  &displayName,
				ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBSOURCEIPPERSISTENCEPROFILE,
				Purge:        strPtr(model.LBSourceIpPersistenceProfile_PURGE_FULL),
				Timeout:      int64Ptr(300),
			}, model.LBSourceIpPersistenceProfileBindingType()
		},
		checkRead: func(t *testing.T, d *schema.ResourceData) {
			assert.Equal(t, "FULL", d.Get("purge"))
			assert.Equal(t, 300, d.Get("timeout"))
		},
		cacheAware:   true,
		cacheTypeKey: resourceTypeLBSourceIpPersistenceProfile,
	})
}
