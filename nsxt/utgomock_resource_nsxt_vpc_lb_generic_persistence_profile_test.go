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

func TestMockResourceNsxtVpcLBGenericPersistenceProfile(t *testing.T) {
	runVpcLbPersistenceResourceSuite(t, vpcLbPersistCase{
		name:        "generic",
		newRes:      resourceNsxtVpcLBGenericPersistenceProfile,
		id:          "vpc-lb-generic-1",
		displayName: "Test VPC LB Generic Persistence Profile",
		extraData: map[string]interface{}{
			"timeout": 600,
		},
		newModel: func(displayName string) (interface{}, bindings.BindingType) {
			return model.LBGenericPersistenceProfile{
				DisplayName:  &displayName,
				ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBGENERICPERSISTENCEPROFILE,
				Timeout:      int64Ptr(600),
			}, model.LBGenericPersistenceProfileBindingType()
		},
		checkRead: func(t *testing.T, d *schema.ResourceData) {
			assert.Equal(t, 600, d.Get("timeout"))
		},
	})
}
