/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNsxtPolicyVlanSegment() *schema.Resource {
	segSchema := getPolicyCommonSegmentSchema()
	delete(segSchema, "overlay_id")
	delete(segSchema, "connectivity_path")

	return &schema.Resource{
		Create: resourceNsxtPolicyVlanSegmentCreate,
		Read:   resourceNsxtPolicyVlanSegmentRead,
		Update: resourceNsxtPolicyVlanSegmentUpdate,
		Delete: resourceNsxtPolicyVlanSegmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: segSchema,
	}
}

func resourceNsxtPolicyVlanSegmentCreate(d *schema.ResourceData, m interface{}) error {
	return nsxtPolicySegmentCreate(d, m, true)
}

func resourceNsxtPolicyVlanSegmentRead(d *schema.ResourceData, m interface{}) error {
	return nsxtPolicySegmentRead(d, m, true)
}

func resourceNsxtPolicyVlanSegmentUpdate(d *schema.ResourceData, m interface{}) error {
	return nsxtPolicySegmentUpdate(d, m, true)
}

func resourceNsxtPolicyVlanSegmentDelete(d *schema.ResourceData, m interface{}) error {
	return nsxtPolicySegmentDelete(d, m)
}
