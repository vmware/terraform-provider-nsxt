/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var GitCommit string

func dataSourceNsxtProviderInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtProviderInfoRead,

		Schema: map[string]*schema.Schema{
			"commit": {
				Type:        schema.TypeString,
				Description: "Latest commit hash",
				Computed:    true,
			},
			"date": {
				Type:        schema.TypeString,
				Description: "Date compiled",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtProviderInfoRead(d *schema.ResourceData, m interface{}) error {
	d.SetId("nsxt")
	d.Set("commit", GitCommit)
	d.Set("date", time.Now().Format(time.Stamp))
	return nil
}
