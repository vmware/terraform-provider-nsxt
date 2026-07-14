// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyTier1Gateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTier1GatewayRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"edge_cluster_path": {
				Type:        schema.TypeString,
				Description: "The path of the edge cluster connected to this Tier1 gateway",
				Optional:    true,
				Computed:    true,
			},
			"context": getContextSchemaWithSpec(utl.SessionContextSpec{IsRequired: false, IsComputed: false, IsVpc: false, AllowDefaultProject: false, FromGlobal: true}),
		},
	}
}

func dataSourceNsxtPolicyTier1GatewayRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	objID := d.Get("id").(string)
	displayName := d.Get("display_name").(string)
	lookupKey := objID
	if lookupKey == "" {
		lookupKey = displayName
	}

	if lookupKey != "" && IsCacheEnabled() {
		val, err := gcache.readCache(lookupKey, "Tier1", d, m, connector)
		if err == nil {
			converter := bindings.NewTypeConverter()
			goVal, convErrs := converter.ConvertToGolang(val.(*data.StructValue), model.Tier1BindingType())
			if len(convErrs) == 0 {
				obj, ok := goVal.(model.Tier1)
				if ok {
					id := lookupKey
					if obj.Id != nil {
						id = *obj.Id
					}
					d.SetId(id)
					d.Set("id", id)
					d.Set("display_name", obj.DisplayName)
					d.Set("description", obj.Description)
					d.Set("path", obj.Path)

					// Single edge cluster is not informative for global manager
					if isPolicyGlobalManager(m) {
						d.Set("edge_cluster_path", "")
						return nil
					}

					err := resourceNsxtPolicyTier1GatewayReadEdgeCluster(getSessionContext(d, m), d, connector)
					if err != nil {
						return fmt.Errorf("failed to get Tier1 %s locale-services: %v", id, err)
					}
					return nil
				}
			}
		}
	}

	obj, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "Tier1", nil)
	if err != nil {
		return err
	}

	// Single edge cluster is not informative for global manager
	if isPolicyGlobalManager(m) {
		d.Set("edge_cluster_path", "")
	} else {
		converter := bindings.NewTypeConverter()
		dataValue, errors := converter.ConvertToGolang(obj, model.Tier1BindingType())
		if len(errors) > 0 {
			return errors[0]
		}
		tier1 := dataValue.(model.Tier1)
		err := resourceNsxtPolicyTier1GatewayReadEdgeCluster(getSessionContext(d, m), d, connector)
		if err != nil {
			return fmt.Errorf("failed to get Tier1 %s locale-services: %v", *tier1.Id, err)
		}
	}
	return nil
}
