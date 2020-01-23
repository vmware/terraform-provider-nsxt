/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNsxtPolicySegment() *schema.Resource {
	segSchema := getPolicyCommonSegmentSchema()
	delete(segSchema, "vlan_ids")

	return &schema.Resource{
		Create: resourceNsxtPolicySegmentCreate,
		Read:   resourceNsxtPolicySegmentRead,
		Update: resourceNsxtPolicySegmentUpdate,
		Delete: resourceNsxtPolicySegmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: segSchema,
	}
}

func resourceNsxtPolicySegmentCreate(d *schema.ResourceData, m interface{}) error {
	return nsxtPolicySegmentCreate(d, m, false)
}

func resourceNsxtPolicySegmentRead(d *schema.ResourceData, m interface{}) error {
	return nsxtPolicySegmentRead(d, m, false)
}

func resourceNsxtPolicySegmentUpdate(d *schema.ResourceData, m interface{}) error {
	return nsxtPolicySegmentUpdate(d, m, false)
}

func resourceNsxtPolicySegmentDelete(d *schema.ResourceData, m interface{}) error {
	return nsxtPolicySegmentDelete(d, m)
}
